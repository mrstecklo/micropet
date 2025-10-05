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
	handler := newHttpHandlerMux(logger)

	handler.ServeHTTP(responseRecorder, request)

	if http.StatusInternalServerError != responseRecorder.Code {
		t.Errorf("incorrect code: expected %d, got %d", http.StatusInternalServerError, responseRecorder.Code)
	}
}
