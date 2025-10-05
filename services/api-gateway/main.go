package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := createLogger()

	serveMux := http.NewServeMux()
	serveMux.Handle("/orders/", httpHandler{logger, handleOrders})

	server := http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		Handler:      serveMux,
	}

	logger.Info("Starting server")
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Error("server", slog.String("error", err.Error()))
	}
	logger.Info("Server closed")
}

func createLogger() *slog.Logger {
	options := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewTextHandler(os.Stdout, options)
	return slog.New(handler)
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
