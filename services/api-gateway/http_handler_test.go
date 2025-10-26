package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrstecklo/micropet/services/mock/mock_http"
	"go.uber.org/mock/gomock"
)

type mockServer struct {
	server      *httptest.Server
	mockHandler *mock_http.MockHandler
}

type httpHandlerFixture struct {
	mux              httpHandlerMux
	responseRecorder *httptest.ResponseRecorder
	mockCtrl         *gomock.Controller
	orders           mockServer
}

func setUpHttpHandlerTest(t *testing.T) httpHandlerFixture {
	logger := createLogger()
	mockCtrl := gomock.NewController(t)
	ordersMockHandler := mock_http.NewMockHandler(mockCtrl)
	ordersServer := httptest.NewServer(ordersMockHandler)
	t.Cleanup(ordersServer.Close)
	mux := newHttpHandlerMux(httpHandlerMuxConfig{
		logger: logger,
		orders: server{ordersServer.URL, ordersServer.Client()},
	})
	return httpHandlerFixture{
		mux:              mux,
		responseRecorder: httptest.NewRecorder(),
		mockCtrl:         mockCtrl,
		orders:           mockServer{ordersServer, ordersMockHandler},
	}
}

func TestHttpHandler_GetOrdersNotAllowed(t *testing.T) {
	f := setUpHttpHandlerTest(t)
	f.orders.mockHandler.EXPECT().
		ServeHTTP(gomock.Any(), gomock.Any()).
		Do(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		})

	request := httptest.NewRequest("GET", "/orders", nil)
	f.mux.ServeHTTP(f.responseRecorder, request)

	if http.StatusMethodNotAllowed != f.responseRecorder.Code {
		t.Errorf("incorrect code: expected %d, got %d", http.StatusMethodNotAllowed, f.responseRecorder.Code)
	}
}

func TestHttpHandler_GetOrdersInternalError(t *testing.T) {
	f := setUpHttpHandlerTest(t)
	f.orders.mockHandler.EXPECT().
		ServeHTTP(gomock.Any(), gomock.Any()).
		Do(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		})

	request := httptest.NewRequest("GET", "/orders", nil)
	f.mux.ServeHTTP(f.responseRecorder, request)

	if http.StatusInternalServerError != f.responseRecorder.Code {
		t.Errorf("incorrect code: expected %d, got %d", http.StatusInternalServerError, f.responseRecorder.Code)
	}
}
