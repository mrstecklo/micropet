package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := createLogger()
	handler := newHttpHandlerMux(httpHandlerMuxConfig{
		logger: logger,
		orders: server{},
	})
	server := http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		Handler:      handler,
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
