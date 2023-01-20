package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/devldavydov/promytheus/tests/data"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testResponse struct {
	code        int
	body        string
	contentType string
}

type testRequest struct {
	method string
	url    string
}

func TestMetricsHandler(t *testing.T) {
	tests := []struct {
		name        string
		req         testRequest
		resp        testResponse
		stgInitFunc func(storage.Storage)
	}{
		/// Update metric tests
		{
			name: "update metric: failed GET request",
			req: testRequest{
				method: http.MethodGet,
				url:    "/update/gauge/metric/1",
			},
			resp: testResponse{
				code:        http.StatusMethodNotAllowed,
				body:        "",
				contentType: "",
			},
		},
		{
			name: "update metric: incorrect URL parts count, #1",
			req: testRequest{
				method: http.MethodPost,
				url:    "/update/gauge/metric1/1/2/3",
			},
			resp: testResponse{
				code:        http.StatusNotFound,
				body:        "404 page not found\n",
				contentType: ContentTypeTextPlain,
			},
		},
		{
			name: "update metric: incorrect URL parts count, #2",
			req: testRequest{
				method: http.MethodPost,
				url:    "/update/gauge/metric1",
			},
			resp: testResponse{
				code:        http.StatusNotFound,
				body:        "404 page not found\n",
				contentType: ContentTypeTextPlain,
			},
		},
		{
			name: "update metric: unknown metric type",
			req: testRequest{
				method: http.MethodPost,
				url:    "/update/fuzz/metric1/1",
			},
			resp: testResponse{
				code:        http.StatusNotImplemented,
				body:        http.StatusText(http.StatusNotImplemented),
				contentType: ContentTypeTextPlain,
			},
		},
		{
			name: "update metric: empty metric name",
			req: testRequest{
				method: http.MethodPost,
				url:    "/update/gauge//1",
			},
			resp: testResponse{
				code:        http.StatusBadRequest,
				body:        http.StatusText(http.StatusBadRequest),
				contentType: ContentTypeTextPlain,
			},
		},
		{
			name: "update metric: incorrect gauge val, #1",
			req: testRequest{
				method: http.MethodPost,
				url:    "/update/gauge/metric1/foobar",
			},
			resp: testResponse{
				code:        http.StatusBadRequest,
				body:        http.StatusText(http.StatusBadRequest),
				contentType: ContentTypeTextPlain,
			},
		},
		{
			name: "update metric: incorrect counter val, #1",
			req: testRequest{
				method: http.MethodPost,
				url:    "/update/counter/metric1/foobar",
			},
			resp: testResponse{
				code:        http.StatusBadRequest,
				body:        http.StatusText(http.StatusBadRequest),
				contentType: ContentTypeTextPlain,
			},
		},
		{
			name: "update metric: incorrect counter val, #2",
			req: testRequest{
				method: http.MethodPost,
				url:    "/update/counter/metric1/1.234",
			},
			resp: testResponse{
				code:        http.StatusBadRequest,
				body:        http.StatusText(http.StatusBadRequest),
				contentType: ContentTypeTextPlain,
			},
		},
		{
			name: "update metric: incorrect counter val, #3",
			req: testRequest{
				method: http.MethodPost,
				url:    "/update/counter/metric1/-1234",
			},
			resp: testResponse{
				code:        http.StatusBadRequest,
				body:        http.StatusText(http.StatusBadRequest),
				contentType: ContentTypeTextPlain,
			},
		},
		{
			name: "update metric: correct gauge",
			req: testRequest{
				method: http.MethodPost,
				url:    "/update/gauge/metric1/1.234",
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        http.StatusText(http.StatusOK),
				contentType: ContentTypeTextPlain,
			},
		},
		{
			name: "update metric: correct counter",
			req: testRequest{
				method: http.MethodPost,
				url:    "/update/counter/metric1/1234",
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        http.StatusText(http.StatusOK),
				contentType: ContentTypeTextPlain,
			},
		},
		/// Get metric test
		{
			name: "get metric: failed POST request",
			req: testRequest{
				method: http.MethodPost,
				url:    "/value/counter/metric1",
			},
			resp: testResponse{
				code:        http.StatusMethodNotAllowed,
				body:        "",
				contentType: "",
			},
		},
		{
			name: "get metric: unknown metric type",
			req: testRequest{
				method: http.MethodGet,
				url:    "/value/fuzzbuzz/metric1",
			},
			resp: testResponse{
				code:        http.StatusNotImplemented,
				body:        http.StatusText(http.StatusNotImplemented),
				contentType: ContentTypeTextPlain,
			},
		},
		{
			name: "get metric: not in storage",
			req: testRequest{
				method: http.MethodGet,
				url:    "/value/counter/metric1",
			},
			resp: testResponse{
				code:        http.StatusNotFound,
				body:        http.StatusText(http.StatusNotFound),
				contentType: ContentTypeTextPlain,
			},
		},
		{
			name: "get metric: counter",
			req: testRequest{
				method: http.MethodGet,
				url:    "/value/counter/metric1",
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        "123",
				contentType: ContentTypeTextPlain,
			},
			stgInitFunc: func(s storage.Storage) {
				s.SetCounterMetric("metric1", 123)
			},
		},
		{
			name: "get metric: gauge",
			req: testRequest{
				method: http.MethodGet,
				url:    "/value/gauge/metric1",
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        "1.230",
				contentType: ContentTypeTextPlain,
			},
			stgInitFunc: func(s storage.Storage) {
				s.SetGaugeMetric("metric1", 1.23)
			},
		},
		{
			name: "get metric: gauge rounded",
			req: testRequest{
				method: http.MethodGet,
				url:    "/value/gauge/metric1",
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        "1.235",
				contentType: ContentTypeTextPlain,
			},
			stgInitFunc: func(s storage.Storage) {
				s.SetGaugeMetric("metric1", 1.23456)
			},
		},
		/// Get all metrics page
		{
			name: "get all metrics page: empty",
			req: testRequest{
				method: http.MethodGet,
				url:    "/",
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        data.AllMetricsEmptyResponse,
				contentType: ContentTypeHTML,
			},
		},
		{
			name: "get all metrics page: with data",
			req: testRequest{
				method: http.MethodGet,
				url:    "/",
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        data.AllMetricsResponseWithData,
				contentType: ContentTypeHTML,
			},
			stgInitFunc: func(s storage.Storage) {
				s.SetGaugeMetric("foo", 1.23456)
				s.SetGaugeMetric("bar", 1.23456)
				s.SetCounterMetric("aaa", 1)
				s.SetCounterMetric("aaa", 1)
				s.SetCounterMetric("zzz", 3)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := storage.NewMemStorage()
			if tt.stgInitFunc != nil {
				tt.stgInitFunc(storage)
			}

			metricsHandler := NewMetricsHandler(storage, logrus.New())
			r := NewRouter(metricsHandler)
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, contentType, body := doTestRequest(t, ts, tt.req.method, tt.req.url)

			assert.Equal(t, tt.resp.code, statusCode)
			assert.Equal(t, tt.resp.body, body)
			assert.Equal(t, tt.resp.contentType, contentType)
		})
	}
}

func doTestRequest(t *testing.T, ts *httptest.Server, method, path string) (int, string, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp.StatusCode, resp.Header.Get("Content-Type"), string(respBody)
}
