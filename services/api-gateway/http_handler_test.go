package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpHandler(t *testing.T) {
	logger := createLogger()
	request := httptest.NewRequest("GET", "/orders", nil)
	responseRecorder := httptest.NewRecorder()
	ordersServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}))
	defer ordersServer.Close()
	handler := newHttpHandlerMux(httpHandlerMuxConfig{
		logger: logger,
		orders: server{ordersServer.URL, ordersServer.Client()},
	})

	handler.ServeHTTP(responseRecorder, request)

	if http.StatusMethodNotAllowed != responseRecorder.Code {
		t.Errorf("incorrect code: expected %d, got %d", http.StatusMethodNotAllowed, responseRecorder.Code)
	}
}
