package handler

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	_http "github.com/devldavydov/promytheus/internal/common/http"
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
	method      string
	url         string
	body        *string
	contentType *string
}

func TestMetricsHandler(t *testing.T) {
	s := func(s string) *string { return &s }

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
				contentType: _http.ContentTypeTextPlain,
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
				contentType: _http.ContentTypeTextPlain,
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
				contentType: _http.ContentTypeTextPlain,
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
				contentType: _http.ContentTypeTextPlain,
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
				contentType: _http.ContentTypeTextPlain,
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
				contentType: _http.ContentTypeTextPlain,
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
				contentType: _http.ContentTypeTextPlain,
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
				contentType: _http.ContentTypeTextPlain,
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
				contentType: _http.ContentTypeTextPlain,
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
				contentType: _http.ContentTypeTextPlain,
			},
		},
		/// Update JSON metrics test
		{
			name: "update JSON metric: failed GET request",
			req: testRequest{
				method:      http.MethodGet,
				url:         "/update/",
				body:        s(`{"id": "foo", "type": "gauge", "value": 123.0}`),
				contentType: s(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusMethodNotAllowed,
				body:        "",
				contentType: "",
			},
		},
		{
			name: "update JSON metric: unknown metric type",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        s(`{"id": "foo", "type": "fuzz", "value": 123.0}`),
				contentType: s(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusNotImplemented,
				body:        http.StatusText(http.StatusNotImplemented),
				contentType: _http.ContentTypeTextPlain,
			},
		},
		{
			name: "update JSON metric: empty metric name",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        s(`{"id": "", "type": "gauge", "value": 123.0}`),
				contentType: s(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusBadRequest,
				body:        http.StatusText(http.StatusBadRequest),
				contentType: _http.ContentTypeTextPlain,
			},
		},
		{
			name: "update JSON metric: incorrect gauge val, #1",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        s(`{"id": "foo", "type": "gauge", "value": "abc"}`),
				contentType: s(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusBadRequest,
				body:        http.StatusText(http.StatusBadRequest),
				contentType: _http.ContentTypeTextPlain,
			},
		},
		{
			name: "update JSON metric: incorrect gauge val, #2",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        s(`{"id": "foo", "type": "gauge", "delta": 123.0}`),
				contentType: s(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusBadRequest,
				body:        http.StatusText(http.StatusBadRequest),
				contentType: _http.ContentTypeTextPlain,
			},
		},
		{
			name: "update JSON metric: incorrect counter val, #1",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        s(`{"id": "foo", "type": "counter", "value": 123}`),
				contentType: s(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusBadRequest,
				body:        http.StatusText(http.StatusBadRequest),
				contentType: _http.ContentTypeTextPlain,
			},
		},
		{
			name: "update JSON metric: incorrect counter val, #2",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        s(`{"id": "foo", "type": "counter", "delta": 123.0}`),
				contentType: s(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusBadRequest,
				body:        http.StatusText(http.StatusBadRequest),
				contentType: _http.ContentTypeTextPlain,
			},
		},
		{
			name: "update JSON metric: incorrect counter val, #3",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        s(`{"id": "foo", "type": "counter", "delta": -123}`),
				contentType: s(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusBadRequest,
				body:        http.StatusText(http.StatusBadRequest),
				contentType: _http.ContentTypeTextPlain,
			},
		},
		{
			name: "update JSON metric: correct gauge",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        s(`{"id": "foo", "type": "gauge", "value": 123.0}`),
				contentType: s(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"foo","type":"gauge","value":123}` + "\n",
				contentType: _http.ContentTypeApplicationJSON,
			},
		},
		{
			name: "update JSON metric: correct counter",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        s(`{"id": "foo", "type": "counter", "delta": 123}`),
				contentType: s(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"foo","type":"counter","delta":123}` + "\n",
				contentType: _http.ContentTypeApplicationJSON,
			},
		},
		{
			name: "update JSON metric: correct counter with update",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        s(`{"id": "foo", "type": "counter", "delta": 123}`),
				contentType: s(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"foo","type":"counter","delta":246}` + "\n",
				contentType: _http.ContentTypeApplicationJSON,
			},
			stgInitFunc: func(s storage.Storage) {
				s.SetCounterMetric("foo", 123)
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
				contentType: _http.ContentTypeTextPlain,
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
				contentType: _http.ContentTypeTextPlain,
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
				contentType: _http.ContentTypeTextPlain,
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
				contentType: _http.ContentTypeTextPlain,
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
				contentType: _http.ContentTypeTextPlain,
			},
			stgInitFunc: func(s storage.Storage) {
				s.SetGaugeMetric("metric1", 1.23456)
			},
		},
		/// Get JSON metric test
		{
			name: "get JSON metric: failed GET request",
			req: testRequest{
				method:      http.MethodGet,
				url:         "/value/",
				body:        s(`{"id": "foo", "type": "gauge"}`),
				contentType: s(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusMethodNotAllowed,
				body:        "",
				contentType: "",
			},
		},
		{
			name: "get JSON metric: unknown metric type",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/value/",
				body:        s(`{"id": "foo", "type": "fuzz"}`),
				contentType: s(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusNotImplemented,
				body:        http.StatusText(http.StatusNotImplemented),
				contentType: _http.ContentTypeTextPlain,
			},
		},
		{
			name: "get JSON metric: not in storage",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/value/",
				body:        s(`{"id": "foo", "type": "counter"}`),
				contentType: s(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusNotFound,
				body:        http.StatusText(http.StatusNotFound),
				contentType: _http.ContentTypeTextPlain,
			},
		},
		{
			name: "get JSON metric: gauge",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/value/",
				body:        s(`{"id": "foo", "type": "gauge"}`),
				contentType: s(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"foo","type":"gauge","value":123}` + "\n",
				contentType: _http.ContentTypeApplicationJSON,
			},
			stgInitFunc: func(s storage.Storage) {
				s.SetGaugeMetric("foo", 123.0)
			},
		},
		{
			name: "get JSON metric: gauge rounded",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/value/",
				body:        s(`{"id": "foo", "type": "gauge"}`),
				contentType: s(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"foo","type":"gauge","value":1.23456}` + "\n",
				contentType: _http.ContentTypeApplicationJSON,
			},
			stgInitFunc: func(s storage.Storage) {
				s.SetGaugeMetric("foo", 1.23456)
			},
		},
		{
			name: "get JSON metric: counter",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/value/",
				body:        s(`{"id": "foo", "type": "counter"}`),
				contentType: s(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"foo","type":"counter","delta":123}` + "\n",
				contentType: _http.ContentTypeApplicationJSON,
			},
			stgInitFunc: func(s storage.Storage) {
				s.SetCounterMetric("foo", 123)
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
				contentType: _http.ContentTypeHTML,
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
				contentType: _http.ContentTypeHTML,
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
			logger := logrus.New()

			storage, _ := storage.NewMemStorage(context.TODO(), logger, storage.NewPersistSettings(0, "", false))
			if tt.stgInitFunc != nil {
				tt.stgInitFunc(storage)
			}

			metricsHandler := NewMetricsHandler(storage, logger)
			r := NewRouter(metricsHandler)
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, contentType, body := doTestRequest(t, ts, tt.req)

			assert.Equal(t, tt.resp.code, statusCode)
			assert.Equal(t, tt.resp.body, body)
			assert.Equal(t, tt.resp.contentType, contentType)
		})
	}
}

func doTestRequest(t *testing.T, ts *httptest.Server, testReq testRequest) (int, string, string) {
	var buf io.Reader
	if testReq.body != nil {
		buf = bytes.NewBuffer([]byte(*testReq.body))
	}

	req, err := http.NewRequest(testReq.method, ts.URL+testReq.url, buf)
	require.NoError(t, err)

	if testReq.contentType != nil {
		req.Header.Set("Content-Type", *testReq.contentType)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp.StatusCode, resp.Header.Get("Content-Type"), string(respBody)
}
