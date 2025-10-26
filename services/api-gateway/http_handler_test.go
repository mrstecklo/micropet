package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/mrstecklo/micropet/services/mock/mock_http"
	"github.com/stretchr/testify/assert"
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
		orders: serverConfig{ordersServer.URL, ordersServer.Client()},
	})
	return httpHandlerFixture{
		mux:              mux,
		responseRecorder: httptest.NewRecorder(),
		mockCtrl:         mockCtrl,
		orders:           mockServer{ordersServer, ordersMockHandler},
	}
}

func TestHttpHandler_ReturnsErrorsFromOrders(t *testing.T) {
	data := []struct {
		name   string
		method string
		target string
		code   int
		body   string
	}{
		{
			"GetRootMethodNotAllowed",
			"GET",
			"/orders",
			http.StatusMethodNotAllowed,
			"Method not allowed",
		},
		{
			"GetRootInternalServerError",
			"GET",
			"/orders",
			http.StatusInternalServerError,
			"Internal server error",
		},
		{
			"GetIdInternalServerError",
			"GET",
			"/orders/123",
			http.StatusInternalServerError,
			"Internal server error",
		},
		{
			"PostIdMethodNotAllowed",
			"POST",
			"/orders/123",
			http.StatusMethodNotAllowed,
			"Method not allowed",
		},
		{
			"PostIdInternalServerError",
			"POST",
			"/orders/123",
			http.StatusInternalServerError,
			"Internal server error",
		},
		{
			"PostRootInternalServerError",
			"POST",
			"/orders/123",
			http.StatusInternalServerError,
			"Internal server error",
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			f := setUpHttpHandlerTest(t)
			f.orders.mockHandler.EXPECT().
				ServeHTTP(gomock.Any(), gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, d.body, d.code)
				})

			request := httptest.NewRequest(d.method, d.target, nil)
			f.mux.ServeHTTP(f.responseRecorder, request)

			assert.Equal(t, d.code, f.responseRecorder.Code)
			assert.Equal(t, d.body+"\n", f.responseRecorder.Body.String())
		})
	}
}

func TestHttpHandler_ForwardsRequestMethodAndPathToOrders(t *testing.T) {
	data := []struct {
		name   string
		method string
		target string
	}{
		{
			"PostRoot",
			"POST",
			"/orders",
		},
		{
			"PostId",
			"POST",
			"/orders/111",
		},
		{
			"GetRoot",
			"GET",
			"/orders",
		},
		{
			"GetId",
			"GET",
			"/orders/111",
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			f := setUpHttpHandlerTest(t)

			f.orders.mockHandler.EXPECT().
				ServeHTTP(gomock.Any(), gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, d.method, r.Method)
					assert.Equal(t, d.target, r.URL.Path)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				})

			request := httptest.NewRequest(d.method, d.target, nil)
			f.mux.ServeHTTP(f.responseRecorder, request)
		})
	}
}

func TestHttpHandler_ForwardsRequestBodyToOrders(t *testing.T) {
	data := []struct {
		body string
	}{
		{
			`{"title": "something"}`,
		},
		{
			`{"title": "ololo"}`,
		},
	}
	for idx, d := range data {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			f := setUpHttpHandlerTest(t)

			f.orders.mockHandler.EXPECT().
				ServeHTTP(gomock.Any(), gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request) {
					bytes, err := io.ReadAll(r.Body)
					assert.Nil(t, err)
					assert.Equal(t, d.body, string(bytes))
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				})

			request := httptest.NewRequest("POST", "/orders", strings.NewReader(d.body))
			f.mux.ServeHTTP(f.responseRecorder, request)
		})
	}
}

func TestHttpHandler_ForwardsRequestHeaderToOrders(t *testing.T) {
	f := setUpHttpHandlerTest(t)

	f.orders.mockHandler.EXPECT().
		ServeHTTP(gomock.Any(), gomock.Any()).
		Do(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "expected value", r.Header.Get("X-Testing"))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		})

	request := httptest.NewRequest("GET", "/orders", nil)
	request.Header.Add("X-Testing", "expected value")
	f.mux.ServeHTTP(f.responseRecorder, request)
}

func TestHttpHandler_ForwardsOtherRequestHeaderToOrders(t *testing.T) {
	f := setUpHttpHandlerTest(t)

	f.orders.mockHandler.EXPECT().
		ServeHTTP(gomock.Any(), gomock.Any()).
		Do(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "something", r.Header.Get("X-Something"))
			assert.Equal(t, "quick brown fox", r.Header.Get("X-Foo"))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		})

	request := httptest.NewRequest("GET", "/orders", nil)
	request.Header.Add("X-Something", "something")
	request.Header.Add("X-Foo", "quick brown fox")
	f.mux.ServeHTTP(f.responseRecorder, request)
}
