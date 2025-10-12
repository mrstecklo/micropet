package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type httpHandlerFixtureConfig struct {
	ordersServerHandler func(http.ResponseWriter, *http.Request)
}

type httpHandlerFixture struct {
	handler          httpHandlerMux
	responseRecorder *httptest.ResponseRecorder
	ordersServer     *httptest.Server
}

func setUpHttpHandler(t *testing.T, config httpHandlerFixtureConfig) httpHandlerFixture {
	logger := createLogger()
	ordersServer := httptest.NewServer(http.HandlerFunc(config.ordersServerHandler))
	handler := newHttpHandlerMux(httpHandlerMuxConfig{
		logger: logger,
		orders: server{ordersServer.URL, ordersServer.Client()},
	})
	t.Cleanup(ordersServer.Close)
	return httpHandlerFixture{
		handler:          handler,
		responseRecorder: httptest.NewRecorder(),
		ordersServer:     ordersServer,
	}
}

func TestHttpHandler_GetOrdersNotAllowed(t *testing.T) {
	config := httpHandlerFixtureConfig{
		ordersServerHandler: func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		},
	}
	f := setUpHttpHandler(t, config)
	request := httptest.NewRequest("GET", "/orders", nil)

	f.handler.ServeHTTP(f.responseRecorder, request)

	if http.StatusMethodNotAllowed != f.responseRecorder.Code {
		t.Errorf("incorrect code: expected %d, got %d", http.StatusMethodNotAllowed, f.responseRecorder.Code)
	}
}

func TestHttpHandler_GetOrdersInternalError(t *testing.T) {
	config := httpHandlerFixtureConfig{
		ordersServerHandler: func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Method not allowed", http.StatusInternalServerError)
		},
	}
	f := setUpHttpHandler(t, config)
	request := httptest.NewRequest("GET", "/orders", nil)

	f.handler.ServeHTTP(f.responseRecorder, request)

	if http.StatusInternalServerError != f.responseRecorder.Code {
		t.Errorf("incorrect code: expected %d, got %d", http.StatusInternalServerError, f.responseRecorder.Code)
	}
}
