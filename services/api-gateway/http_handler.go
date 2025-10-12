package main

import (
	"log/slog"
	"net/http"
)

type server struct {
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
	orders server
}

func newHttpHandlerMux(config httpHandlerMuxConfig) httpHandlerMux {
	mux := http.NewServeMux()
	mux.Handle("/orders", httpHandler{config.logger, handleOrders, config.orders})
	return httpHandlerMux{mux}
}

type handleFunc func(*slog.Logger, server, http.ResponseWriter, *http.Request)

type httpHandler struct {
	logger *slog.Logger
	handle handleFunc
	server server
}

func (h httpHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	h.handle(h.logger, h.server, responseWriter, request)
}

func handleOrders(logger *slog.Logger, server server, responseWriter http.ResponseWriter, request *http.Request) {
	req, _ := http.NewRequest("GET", server.url, nil)
	resp, _ := server.client.Do(req)
	responseWriter.WriteHeader(resp.StatusCode)
}
