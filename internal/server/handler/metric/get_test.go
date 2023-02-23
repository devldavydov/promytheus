package metric

import (
	"errors"
	"net/http"
	"testing"

	_http "github.com/devldavydov/promytheus/internal/common/http"
	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/devldavydov/promytheus/internal/server/mocks"
	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/devldavydov/promytheus/tests/data"
)

func TestGetMetric(t *testing.T) {
	tests := []testItem{
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
	}

	runTests(t, tests)
}

func TestGetMetricFromDb(t *testing.T) {
	tests := []testItem{
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
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().GetCounterMetric("metric1").Return(metric.Counter(123), nil)
			},
		},
		{
			name: "get metric: counter not found",
			req: testRequest{
				method: http.MethodGet,
				url:    "/value/counter/metric1",
			},
			resp: testResponse{
				code:        http.StatusNotFound,
				body:        http.StatusText(http.StatusNotFound),
				contentType: _http.ContentTypeTextPlain,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().GetCounterMetric("metric1").Return(metric.Counter(0), storage.ErrMetricNotFound)
			},
		},
		{
			name: "get metric: counter db err",
			req: testRequest{
				method: http.MethodGet,
				url:    "/value/counter/metric1",
			},
			resp: testResponse{
				code:        http.StatusInternalServerError,
				body:        http.StatusText(http.StatusInternalServerError),
				contentType: _http.ContentTypeTextPlain,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().GetCounterMetric("metric1").Return(metric.Counter(0), errors.New("db error"))
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
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().GetGaugeMetric("metric1").Return(metric.Gauge(1.230), nil)
			},
		},
		{
			name: "get metric: gauge not found",
			req: testRequest{
				method: http.MethodGet,
				url:    "/value/gauge/metric1",
			},
			resp: testResponse{
				code:        http.StatusNotFound,
				body:        http.StatusText(http.StatusNotFound),
				contentType: _http.ContentTypeTextPlain,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().GetGaugeMetric("metric1").Return(metric.Gauge(0), storage.ErrMetricNotFound)
			},
		},
		{
			name: "get metric: gauge db err",
			req: testRequest{
				method: http.MethodGet,
				url:    "/value/gauge/metric1",
			},
			resp: testResponse{
				code:        http.StatusInternalServerError,
				body:        http.StatusText(http.StatusInternalServerError),
				contentType: _http.ContentTypeTextPlain,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().GetGaugeMetric("metric1").Return(metric.Gauge(0), errors.New("db error"))
			},
		},
	}

	runTests(t, tests)
}

func TestGetJSONMetric(t *testing.T) {
	tests := []testItem{
		/// Get JSON Metric
		{
			name: "get JSON metric: failed GET request",
			req: testRequest{
				method:      http.MethodGet,
				url:         "/value/",
				body:        bodyStringReader(`{"id": "foo", "type": "gauge"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
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
				body:        bodyStringReader(`{"id": "foo", "type": "fuzz"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
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
				body:        bodyStringReader(`{"id": "foo", "type": "counter"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
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
				body:        bodyStringReader(`{"id": "foo", "type": "gauge"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"foo","type":"gauge","value":123}`,
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
				body:        bodyStringReader(`{"id": "foo", "type": "gauge"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"foo","type":"gauge","value":1.23456}`,
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
				body:        bodyStringReader(`{"id": "foo", "type": "counter"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"foo","type":"counter","delta":123}`,
				contentType: _http.ContentTypeApplicationJSON,
			},
			stgInitFunc: func(s storage.Storage) {
				s.SetCounterMetric("foo", 123)
			},
		},
	}

	runTests(t, tests)
}

func TestGetJSONMetricFromDb(t *testing.T) {
	tests := []testItem{
		{
			name: "get JSON metric: gauge",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/value/",
				body:        bodyStringReader(`{"id": "foo", "type": "gauge"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"foo","type":"gauge","value":1.23}`,
				contentType: _http.ContentTypeApplicationJSON,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().GetGaugeMetric("foo").Return(metric.Gauge(1.230), nil)
			},
		},
		{
			name: "get JSON metric: gauge not found",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/value/",
				body:        bodyStringReader(`{"id": "foo", "type": "gauge"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusNotFound,
				body:        http.StatusText(http.StatusNotFound),
				contentType: _http.ContentTypeTextPlain,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().GetGaugeMetric("foo").Return(metric.Gauge(0), storage.ErrMetricNotFound)
			},
		},
		{
			name: "get JSON metric: gauge db error",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/value/",
				body:        bodyStringReader(`{"id": "foo", "type": "gauge"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusInternalServerError,
				body:        http.StatusText(http.StatusInternalServerError),
				contentType: _http.ContentTypeTextPlain,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().GetGaugeMetric("foo").Return(metric.Gauge(0), errors.New("db error"))
			},
		},
		{
			name: "get JSON metric: counter",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/value/",
				body:        bodyStringReader(`{"id": "foo", "type": "counter"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"foo","type":"counter","delta":123}`,
				contentType: _http.ContentTypeApplicationJSON,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().GetCounterMetric("foo").Return(metric.Counter(123), nil)
			},
		},
		{
			name: "get JSON metric: counter not found",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/value/",
				body:        bodyStringReader(`{"id": "foo", "type": "counter"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusNotFound,
				body:        http.StatusText(http.StatusNotFound),
				contentType: _http.ContentTypeTextPlain,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().GetCounterMetric("foo").Return(metric.Counter(0), storage.ErrMetricNotFound)
			},
		},
		{
			name: "get JSON metric: counter db error",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/value/",
				body:        bodyStringReader(`{"id": "foo", "type": "counter"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusInternalServerError,
				body:        http.StatusText(http.StatusInternalServerError),
				contentType: _http.ContentTypeTextPlain,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().GetCounterMetric("foo").Return(metric.Counter(0), errors.New("db error"))
			},
		},
	}

	runTests(t, tests)
}

func TestGetJSONMetricWithHash(t *testing.T) {
	tests := []testItem{
		{
			name: "get JSON metric with hash: gauge",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/value/",
				body:        bodyStringReader(`{"id":"Sys", "type":"gauge"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				hmacKey:     strPointer("foobar"),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"Sys","type":"gauge","value":13220880,"hash":"48a93e5dde0297029bf66cc10a1cdda9be6f858667ea885dc1b0d810032aa292"}`,
				contentType: _http.ContentTypeApplicationJSON,
			},
			stgInitFunc: func(s storage.Storage) {
				s.SetGaugeMetric("Sys", 13220880)
			},
		},
		{
			name: "get JSON metric with hash: counter",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/value/",
				body:        bodyStringReader(`{"id":"PollCount", "type":"counter"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				hmacKey:     strPointer("foobar"),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"PollCount","type":"counter","delta":5,"hash":"b9203cac5904e73da2504aabfb77a419d3d3f9a0baee3707c55070432c6ff5a8"}`,
				contentType: _http.ContentTypeApplicationJSON,
			},
			stgInitFunc: func(s storage.Storage) {
				s.SetCounterMetric("PollCount", 5)
			},
		},
	}

	runTests(t, tests)
}

func TestGetJSONMetricWithHashFromDb(t *testing.T) {
	tests := []testItem{
		{
			name: "get JSON metric with hash: gauge",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/value/",
				body:        bodyStringReader(`{"id":"Sys", "type":"gauge"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				hmacKey:     strPointer("foobar"),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"Sys","type":"gauge","value":13220880,"hash":"48a93e5dde0297029bf66cc10a1cdda9be6f858667ea885dc1b0d810032aa292"}`,
				contentType: _http.ContentTypeApplicationJSON,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().GetGaugeMetric("Sys").Return(metric.Gauge(13220880), nil)
			},
		},
		{
			name: "get JSON metric with hash: counter",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/value/",
				body:        bodyStringReader(`{"id":"PollCount", "type":"counter"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				hmacKey:     strPointer("foobar"),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"PollCount","type":"counter","delta":5,"hash":"b9203cac5904e73da2504aabfb77a419d3d3f9a0baee3707c55070432c6ff5a8"}`,
				contentType: _http.ContentTypeApplicationJSON,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().GetCounterMetric("PollCount").Return(metric.Counter(5), nil)
			},
		},
	}

	runTests(t, tests)
}

func TestGetAllMetricsPage(t *testing.T) {
	tests := []testItem{
		{
			name: "get all metrics page: failed POST request",
			req: testRequest{
				method: http.MethodPost,
				url:    "/",
			},
			resp: testResponse{
				code:        http.StatusMethodNotAllowed,
				body:        "",
				contentType: "",
			},
		},
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
			name: "get all metrics page: empty, gzipped",
			req: testRequest{
				method:  http.MethodGet,
				url:     "/",
				headers: map[string][]string{"Accept-Encoding": {"gzip"}},
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        data.AllMetricsEmptyResponse,
				contentType: _http.ContentTypeHTML,
				headers:     map[string][]string{"Content-Encoding": {"gzip"}},
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

	runTests(t, tests)
}

func TestGetAllMetricsPageFromDb(t *testing.T) {
	tests := []testItem{
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
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().GetAllMetrics().Return([]storage.StorageItem{
					{MetricName: "aaa", Value: metric.Counter(2)},
					{MetricName: "zzz", Value: metric.Counter(3)},
					{MetricName: "bar", Value: metric.Gauge(1.235)},
					{MetricName: "foo", Value: metric.Gauge(1.235)},
				}, nil)
			},
		},
		{
			name: "get all metrics page: db error",
			req: testRequest{
				method: http.MethodGet,
				url:    "/",
			},
			resp: testResponse{
				code:        http.StatusInternalServerError,
				body:        http.StatusText(http.StatusInternalServerError),
				contentType: _http.ContentTypeTextPlain,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().GetAllMetrics().Return(nil, errors.New("db error"))
			},
		},
	}

	runTests(t, tests)
}
