package main

import (
	"log/slog"
	"net/http"
)

type httpHandlerMux struct {
	mux *http.ServeMux
}

func (h httpHandlerMux) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	h.mux.ServeHTTP(responseWriter, request)
}

func newHttpHandlerMux(logger *slog.Logger) httpHandlerMux {
	mux := http.NewServeMux()
	mux.Handle("/orders", httpHandler{logger, handleOrders})
	return httpHandlerMux{mux}
}

type handleFunc func(*slog.Logger, http.ResponseWriter, *http.Request)

type httpHandler struct {
	logger *slog.Logger
	handle handleFunc
}

func (h httpHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	h.handle(h.logger, responseWriter, request)
}

func handleOrders(logger *slog.Logger, responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.WriteHeader(http.StatusInternalServerError)
}
