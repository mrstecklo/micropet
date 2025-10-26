package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrstecklo/micropet/services/mock/mock_http"
	"go.uber.org/mock/gomock"
)

type httpHandlerFixtureConfig struct {
	ordersServerHandler http.Handler
}

type httpHandlerFixture struct {
	handler          httpHandlerMux
	responseRecorder *httptest.ResponseRecorder
	ordersServer     *httptest.Server
}

func setUpHttpHandler(t *testing.T, config httpHandlerFixtureConfig) httpHandlerFixture {
	logger := createLogger()
	ordersServer := httptest.NewServer(config.ordersServerHandler)
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
	ctrl := gomock.NewController(t)
	ordersHandlerMock := mock_http.NewMockHandler(ctrl)
	f := setUpHttpHandler(t, httpHandlerFixtureConfig{
		ordersServerHandler: ordersHandlerMock,
	})
	request := httptest.NewRequest("GET", "/orders", nil)

	ordersHandlerMock.EXPECT().
		ServeHTTP(gomock.Any(), gomock.Any()).
		Do(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		})

	f.handler.ServeHTTP(f.responseRecorder, request)

	if http.StatusMethodNotAllowed != f.responseRecorder.Code {
		t.Errorf("incorrect code: expected %d, got %d", http.StatusMethodNotAllowed, f.responseRecorder.Code)
	}
}

func TestHttpHandler_GetOrdersInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	ordersHandlerMock := mock_http.NewMockHandler(ctrl)
	f := setUpHttpHandler(t, httpHandlerFixtureConfig{
		ordersServerHandler: ordersHandlerMock,
	})
	request := httptest.NewRequest("GET", "/orders", nil)

	ordersHandlerMock.EXPECT().
		ServeHTTP(gomock.Any(), gomock.Any()).
		Do(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		})

	f.handler.ServeHTTP(f.responseRecorder, request)

	if http.StatusInternalServerError != f.responseRecorder.Code {
		t.Errorf("incorrect code: expected %d, got %d", http.StatusInternalServerError, f.responseRecorder.Code)
	}
}
