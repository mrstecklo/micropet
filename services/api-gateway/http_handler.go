package main

import (
	"io"
	"log/slog"
	"net/http"
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
	req, _ := http.NewRequest("GET", server.url, nil)
	resp, _ := server.client.Do(req)
	responseWriter.WriteHeader(resp.StatusCode)
	io.Copy(responseWriter, resp.Body)
}
