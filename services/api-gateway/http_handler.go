package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

type serverConfig struct {
	url    string
	client *http.Client
}

type httpHandlerMux struct {
	mux *http.ServeMux
}

func (h httpHandlerMux) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	h.mux.ServeHTTP(responseWriter, request)
}

type httpHandlerMuxConfig struct {
	logger *slog.Logger
	orders serverConfig
}

func newHttpHandlerMux(config httpHandlerMuxConfig) httpHandlerMux {
	mux := http.NewServeMux()
	ordersHandler := httpHandler{config.logger, handleOrders, config.orders}
	mux.Handle("/orders", ordersHandler)
	mux.Handle("/orders/", ordersHandler)
	return httpHandlerMux{mux}
}

type handleFunc func(*slog.Logger, serverConfig, http.ResponseWriter, *http.Request)

type httpHandler struct {
	logger *slog.Logger
	handle handleFunc
	server serverConfig
}

func (h httpHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	h.handle(h.logger, h.server, responseWriter, request)
}

func handleOrders(logger *slog.Logger, server serverConfig, responseWriter http.ResponseWriter, request *http.Request) {
	logger.Info("handle orders", "method", request.Method, "url", request.URL.String())
	proxyURL, err := url.Parse(server.url)
	if err != nil {
		logger.Error("failed to parse url", "error", err.Error(), "url", server.url)
		http.Error(responseWriter, "Internal server error", http.StatusInternalServerError)
		return
	}
	proxyURL.Path = request.URL.Path
	proxyRequest, err := http.NewRequest(request.Method, proxyURL.String(), request.Body)
	if err != nil {
		logger.Error("failed to create http request", "error", err.Error(), "method", request.Method, "url", proxyURL.String())
		http.Error(responseWriter, "Internal server error", http.StatusInternalServerError)
		return
	}
	proxyRequest.Header = request.Header.Clone()
	proxyResponse, err := server.client.Do(proxyRequest)
	if err != nil {
		logger.Error("failed to send http request", "error", err.Error(), "method", request.Method, "url", proxyURL.String())
		http.Error(responseWriter, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer func() {
		err := proxyResponse.Body.Close()
		if err != nil {
			logger.Error("failed to close response body", "error", err.Error())
		}
	}()

	for key, value := range proxyResponse.Header {
		responseWriter.Header()[key] = value
	}
	responseWriter.WriteHeader(proxyResponse.StatusCode)
	_, err = io.Copy(responseWriter, proxyResponse.Body)
	if err != nil {
		logger.Error("failed to copy response body", "error", err.Error())
	}
}
