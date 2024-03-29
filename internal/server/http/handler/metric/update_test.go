package metric

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	_http "github.com/devldavydov/promytheus/internal/common/http"
	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/devldavydov/promytheus/internal/server/mocks"
	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/go-resty/resty/v2"
)

func TestUpdateMetric(t *testing.T) {
	tests := []testItem{
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
		{
			name: "update metric: correct gauge, wrong subnet",
			req: testRequest{
				method:  http.MethodPost,
				url:     "/update/gauge/metric2/1.234",
				headers: map[string][]string{"X-Real-IP": {"1.1.1.1"}},
			},
			resp: testResponse{
				code:        http.StatusForbidden,
				body:        "",
				contentType: "",
			},
			trustedSubnet: getIPNet("10.0.0.0/16"),
		},
		{
			name: "update metric: correct counter, wrong subnet",
			req: testRequest{
				method:  http.MethodPost,
				url:     "/update/counter/metric2/1234",
				headers: map[string][]string{"X-Real-IP": {"1.1.1.1"}},
			},
			resp: testResponse{
				code:        http.StatusForbidden,
				body:        "",
				contentType: "",
			},
			trustedSubnet: getIPNet("10.0.0.0/16"),
		},
		{
			name: "update metric: correct gauge, correct subnet",
			req: testRequest{
				method:  http.MethodPost,
				url:     "/update/gauge/metric2/1.234",
				headers: map[string][]string{"X-Real-IP": {"192.168.0.100"}},
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        http.StatusText(http.StatusOK),
				contentType: _http.ContentTypeTextPlain,
			},
			trustedSubnet: getIPNet("192.168.0.0/16"),
		},
		{
			name: "update metric: correct counter, correct subnet",
			req: testRequest{
				method:  http.MethodPost,
				url:     "/update/counter/metric2/1234",
				headers: map[string][]string{"X-Real-IP": {"192.168.0.100"}},
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        http.StatusText(http.StatusOK),
				contentType: _http.ContentTypeTextPlain,
			},
			trustedSubnet: getIPNet("192.168.0.0/16"),
		},
	}

	runTests(t, tests)
}

func TestUpdateMetricInDb(t *testing.T) {
	tests := []testItem{
		{
			name: "update metric: gauge",
			req: testRequest{
				method: http.MethodPost,
				url:    "/update/gauge/metric1/1.234",
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        http.StatusText(http.StatusOK),
				contentType: _http.ContentTypeTextPlain,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().SetGaugeMetric("metric1", metric.Gauge(1.234)).Return(metric.Gauge(1.234), nil)
			},
		},
		{
			name: "update metric: gauge db err",
			req: testRequest{
				method: http.MethodPost,
				url:    "/update/gauge/metric1/1.234",
			},
			resp: testResponse{
				code:        http.StatusInternalServerError,
				body:        http.StatusText(http.StatusInternalServerError),
				contentType: _http.ContentTypeTextPlain,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().SetGaugeMetric("metric1", metric.Gauge(1.234)).Return(metric.Gauge(0), errors.New("db error"))
			},
		},
		{
			name: "update metric: counter",
			req: testRequest{
				method: http.MethodPost,
				url:    "/update/counter/metric1/1234",
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        http.StatusText(http.StatusOK),
				contentType: _http.ContentTypeTextPlain,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().SetCounterMetric("metric1", metric.Counter(1234)).Return(metric.Counter(1234), nil)
			},
		},
		{
			name: "update metric: counter db err",
			req: testRequest{
				method: http.MethodPost,
				url:    "/update/counter/metric1/1234",
			},
			resp: testResponse{
				code:        http.StatusInternalServerError,
				body:        http.StatusText(http.StatusInternalServerError),
				contentType: _http.ContentTypeTextPlain,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().SetCounterMetric("metric1", metric.Counter(1234)).Return(metric.Counter(0), errors.New("db error"))
			},
		},
	}

	runTests(t, tests)
}

func TestUpdateJsonMetric(t *testing.T) {
	tests := []testItem{
		{
			name: "update JSON metric: failed GET request",
			req: testRequest{
				method:      http.MethodGet,
				url:         "/update/",
				body:        bodyStringReader(`{"id": "foo", "type": "gauge", "value": 123.0}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
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
				body:        bodyStringReader(`{"id": "foo", "type": "fuzz", "value": 123.0}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
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
				body:        bodyStringReader(`{"id": "", "type": "gauge", "value": 123.0}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
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
				body:        bodyStringReader(`{"id": "foo", "type": "gauge", "value": "abc"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
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
				body:        bodyStringReader(`{"id": "foo", "type": "gauge", "delta": 123.0}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
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
				body:        bodyStringReader(`{"id": "foo", "type": "counter", "value": 123}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
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
				body:        bodyStringReader(`{"id": "foo", "type": "counter", "delta": 123.0}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
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
				body:        bodyStringReader(`{"id": "foo", "type": "counter", "delta": -123}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
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
				body:        bodyStringReader(`{"id": "foo", "type": "gauge", "value": 123.0}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"foo","type":"gauge","value":123}`,
				contentType: _http.ContentTypeApplicationJSON,
			},
		},
		{
			name: "update JSON metric: correct counter",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        bodyStringReader(`{"id": "foo", "type": "counter", "delta": 123}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"foo","type":"counter","delta":123}`,
				contentType: _http.ContentTypeApplicationJSON,
			},
		},
		{
			name: "update JSON metric: correct counter with update",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        bodyStringReader(`{"id": "foo", "type": "counter", "delta": 123}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"foo","type":"counter","delta":246}`,
				contentType: _http.ContentTypeApplicationJSON,
			},
			stgInitFunc: func(s storage.Storage) {
				s.SetCounterMetric("foo", 123)
			},
		},
		{
			name: "update JSON metric: correct counter with update, gzipped",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        bodyGzipReader(`{"id": "foo", "type": "counter", "delta": 123}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				headers:     map[string][]string{"Accept-Encoding": {"gzip"}, "Content-Encoding": {"gzip"}},
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"foo","type":"counter","delta":246}`,
				contentType: _http.ContentTypeApplicationJSON,
				headers:     map[string][]string{"Content-Encoding": {"gzip"}},
			},
			stgInitFunc: func(s storage.Storage) {
				s.SetCounterMetric("foo", 123)
			},
		},
		{
			name: "update JSON metric: correct gauge, encrypted",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        bodyStringReader(encryptString(`{"id": "foo_encr", "type": "gauge", "value": 123.0}`)),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				encryption:  true,
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"foo_encr","type":"gauge","value":123}`,
				contentType: _http.ContentTypeApplicationJSON,
			},
		},
		{
			name: "update JSON metric: correct counter, encrypted",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        bodyStringReader(encryptString(`{"id": "foo_encr", "type": "counter", "delta": 123}`)),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				encryption:  true,
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"foo_encr","type":"counter","delta":123}`,
				contentType: _http.ContentTypeApplicationJSON,
			},
		},
		{
			name: "update JSON metric: correct counter with encryption, gzipped",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        bodyGzipReader(encryptString(`{"id": "foo_encr_gz", "type": "counter", "delta": 123}`)),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				headers:     map[string][]string{"Accept-Encoding": {"gzip"}, "Content-Encoding": {"gzip"}},
				encryption:  true,
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"foo_encr_gz","type":"counter","delta":123}`,
				contentType: _http.ContentTypeApplicationJSON,
				headers:     map[string][]string{"Content-Encoding": {"gzip"}},
			},
			stgInitFunc: func(s storage.Storage) {
				s.SetCounterMetric("foo", 123)
			},
		},
		{
			name: "update JSON metric: correct counter with encryption, gzipped, with trusted network",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        bodyGzipReader(encryptString(`{"id": "foo_encr_gz2", "type": "counter", "delta": 123}`)),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				headers:     map[string][]string{"Accept-Encoding": {"gzip"}, "Content-Encoding": {"gzip"}, "X-Real-IP": {"192.168.0.100"}},
				encryption:  true,
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"foo_encr_gz2","type":"counter","delta":246}`,
				contentType: _http.ContentTypeApplicationJSON,
				headers:     map[string][]string{"Content-Encoding": {"gzip"}},
			},
			trustedSubnet: getIPNet("192.168.0.0/16"),
			stgInitFunc: func(s storage.Storage) {
				s.SetCounterMetric("foo_encr_gz2", 123)
			},
		},
		{
			name: "update JSON metric: correct counter with encryption, gzipped, with wrong network",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        bodyGzipReader(encryptString(`{"id": "foo_encr_gz3", "type": "counter", "delta": 123}`)),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				headers:     map[string][]string{"Accept-Encoding": {"gzip"}, "Content-Encoding": {"gzip"}, "X-Real-IP": {"10.10.10.10"}},
				encryption:  true,
			},
			resp: testResponse{
				code:        http.StatusForbidden,
				body:        "",
				contentType: "",
				headers:     map[string][]string{"Content-Encoding": {"gzip"}},
			},
			trustedSubnet: getIPNet("192.168.0.0/16"),
			stgInitFunc: func(s storage.Storage) {
				s.SetCounterMetric("foo_encr_gz3", 123)
			},
		},
	}

	runTests(t, tests)
}

func TestUpdateJsonMetricInDb(t *testing.T) {
	tests := []testItem{
		{
			name: "update JSON metric: gauge",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        bodyStringReader(`{"id": "foo", "type": "gauge", "value": 123.0}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"foo","type":"gauge","value":123}`,
				contentType: _http.ContentTypeApplicationJSON,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().SetGaugeMetric("foo", metric.Gauge(123.0)).Return(metric.Gauge(123.0), nil)
			},
		},
		{
			name: "update JSON metric: gauge db err",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        bodyStringReader(`{"id": "foo", "type": "gauge", "value": 123.0}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusInternalServerError,
				body:        http.StatusText(http.StatusInternalServerError),
				contentType: _http.ContentTypeTextPlain,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().SetGaugeMetric("foo", metric.Gauge(123.0)).Return(metric.Gauge(0), errors.New("db error"))
			},
		},
		{
			name: "update JSON metric: counter",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        bodyStringReader(`{"id": "foo", "type": "counter", "delta": 123}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"foo","type":"counter","delta":123}`,
				contentType: _http.ContentTypeApplicationJSON,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().SetCounterMetric("foo", metric.Counter(123)).Return(metric.Counter(123), nil)
			},
		},
		{
			name: "update JSON metric: counter db err",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        bodyStringReader(`{"id": "foo", "type": "counter", "delta": 123}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusInternalServerError,
				body:        http.StatusText(http.StatusInternalServerError),
				contentType: _http.ContentTypeTextPlain,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().SetCounterMetric("foo", metric.Counter(123)).Return(metric.Counter(0), errors.New("db error"))
			},
		},
	}

	runTests(t, tests)
}

func TestUpdateJSONMetricWithHash(t *testing.T) {
	tests := []testItem{
		{
			name: "update JSON metric with valid hash: gauge",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        bodyStringReader(`{"id":"Sys","type":"gauge","value":13220880,"hash":"48a93e5dde0297029bf66cc10a1cdda9be6f858667ea885dc1b0d810032aa292"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				hmacKey:     strPointer("foobar"),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"Sys","type":"gauge","value":13220880,"hash":"48a93e5dde0297029bf66cc10a1cdda9be6f858667ea885dc1b0d810032aa292"}`,
				contentType: _http.ContentTypeApplicationJSON,
			},
		},
		{
			name: "update JSON metric with invalid hash: gauge",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        bodyStringReader(`{"id":"Sys","type":"gauge","value":13220880,"hash":"48a93e5dde0297029bf66cc10a1cdda9be6f858667ea885dc1b0d810032aa292"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				hmacKey:     strPointer("foobar2"),
			},
			resp: testResponse{
				code:        http.StatusBadRequest,
				body:        http.StatusText(http.StatusBadRequest),
				contentType: _http.ContentTypeTextPlain,
			},
		},
		{
			name: "update JSON metric with valid hash: counter",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        bodyStringReader(`{"id":"PollCount","type":"counter","delta":5,"hash":"b9203cac5904e73da2504aabfb77a419d3d3f9a0baee3707c55070432c6ff5a8"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				hmacKey:     strPointer("foobar"),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `{"id":"PollCount","type":"counter","delta":5,"hash":"b9203cac5904e73da2504aabfb77a419d3d3f9a0baee3707c55070432c6ff5a8"}`,
				contentType: _http.ContentTypeApplicationJSON,
			},
		},
		{
			name: "update JSON metric with invalid hash: counter",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        bodyStringReader(`{"id":"PollCount","type":"counter","delta":5,"hash":"b9203cac5904e73da2504aabfb77a419d3d3f9a0baee3707c55070432c6ff5a8"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				hmacKey:     strPointer("foobar2"),
			},
			resp: testResponse{
				code:        http.StatusBadRequest,
				body:        http.StatusText(http.StatusBadRequest),
				contentType: _http.ContentTypeTextPlain,
			},
		},
	}

	runTests(t, tests)
}

func TestUpdateJSONMetricInDbWithHash(t *testing.T) {
	tests := []testItem{
		{
			name: "update JSON metric with valid hash: gauge",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        bodyStringReader(`{"id":"Sys","type":"gauge","value":13220880,"hash":"48a93e5dde0297029bf66cc10a1cdda9be6f858667ea885dc1b0d810032aa292"}`),
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
				ms.EXPECT().SetGaugeMetric("Sys", metric.Gauge(13220880)).Return(metric.Gauge(13220880), nil)
			},
		},
		{
			name: "update JSON metric with valid hash: gauge db err",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        bodyStringReader(`{"id":"Sys","type":"gauge","value":13220880,"hash":"48a93e5dde0297029bf66cc10a1cdda9be6f858667ea885dc1b0d810032aa292"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				hmacKey:     strPointer("foobar"),
			},
			resp: testResponse{
				code:        http.StatusInternalServerError,
				body:        http.StatusText(http.StatusInternalServerError),
				contentType: _http.ContentTypeTextPlain,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().SetGaugeMetric("Sys", metric.Gauge(13220880)).Return(metric.Gauge(0), errors.New("db error"))
			},
		},
		{
			name: "update JSON metric with valid hash: counter",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        bodyStringReader(`{"id":"PollCount","type":"counter","delta":5,"hash":"b9203cac5904e73da2504aabfb77a419d3d3f9a0baee3707c55070432c6ff5a8"}`),
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
				ms.EXPECT().SetCounterMetric("PollCount", metric.Counter(5)).Return(metric.Counter(5), nil)
			},
		},
		{
			name: "update JSON metric with valid hash: counter db error",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/update/",
				body:        bodyStringReader(`{"id":"PollCount","type":"counter","delta":5,"hash":"b9203cac5904e73da2504aabfb77a419d3d3f9a0baee3707c55070432c6ff5a8"}`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				hmacKey:     strPointer("foobar"),
			},
			resp: testResponse{
				code:        http.StatusInternalServerError,
				body:        http.StatusText(http.StatusInternalServerError),
				contentType: _http.ContentTypeTextPlain,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().SetCounterMetric("PollCount", metric.Counter(5)).Return(metric.Counter(5), errors.New("db error"))
			},
		},
	}

	runTests(t, tests)
}

func TestUpdateMetricJSONBatch(t *testing.T) {
	tests := []testItem{
		{
			name: "update JSON metric: failed GET request",
			req: testRequest{
				method:      http.MethodGet,
				url:         "/updates/",
				body:        bodyStringReader(`[{"id": "foo", "type": "gauge", "value": 123.0}]`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
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
				url:         "/updates/",
				body:        bodyStringReader(`[{"id": "foo", "type": "fuzz", "value": 123.0}]`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
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
				url:         "/updates/",
				body:        bodyStringReader(`[{"id": "", "type": "gauge", "value": 123.0}]`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
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
				url:         "/updates/",
				body:        bodyStringReader(`[{"id": "foo", "type": "gauge", "value": "abc"}]`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
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
				url:         "/updates/",
				body:        bodyStringReader(`[{"id": "foo", "type": "gauge", "delta": 123.0}]`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
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
				url:         "/updates/",
				body:        bodyStringReader(`[{"id": "foo", "type": "counter", "value": 123}]`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
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
				url:         "/updates/",
				body:        bodyStringReader(`[{"id": "foo", "type": "counter", "delta": 123.0}]`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
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
				url:         "/updates/",
				body:        bodyStringReader(`[{"id": "foo", "type": "counter", "delta": -123}]`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
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
				url:         "/updates/",
				body:        bodyStringReader(`[{"id": "foo", "type": "gauge", "value": 123.0}]`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `[]`,
				contentType: _http.ContentTypeApplicationJSON,
			},
			stgCheckFunc: func() []storage.StorageItem {
				return []storage.StorageItem{
					{MetricName: "foo", Value: metric.Gauge(123.0)},
				}
			},
		},
		{
			name: "update JSON metric: correct counter",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/updates/",
				body:        bodyStringReader(`[{"id": "foo", "type": "counter", "delta": 123}]`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `[]`,
				contentType: _http.ContentTypeApplicationJSON,
			},
			stgCheckFunc: func() []storage.StorageItem {
				return []storage.StorageItem{
					{MetricName: "foo", Value: metric.Counter(123)},
				}
			},
		},
		{
			name: "update JSON metric: correct counter with update",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/updates/",
				body:        bodyStringReader(`[{"id": "foo", "type": "counter", "delta": 123}]`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `[]`,
				contentType: _http.ContentTypeApplicationJSON,
			},
			stgInitFunc: func(s storage.Storage) {
				s.SetCounterMetric("foo", 123)
			},
			stgCheckFunc: func() []storage.StorageItem {
				return []storage.StorageItem{
					{MetricName: "foo", Value: metric.Counter(246)},
				}
			},
		},
		{
			name: "update JSON metric: correct counter with update, gzipped",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/updates/",
				body:        bodyGzipReader(`[{"id": "foo", "type": "counter", "delta": 123}]`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				headers:     map[string][]string{"Accept-Encoding": {"gzip"}, "Content-Encoding": {"gzip"}},
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `[]`,
				contentType: _http.ContentTypeApplicationJSON,
				headers:     map[string][]string{"Content-Encoding": {"gzip"}},
			},
			stgInitFunc: func(s storage.Storage) {
				s.SetCounterMetric("foo", 123)
			},
			stgCheckFunc: func() []storage.StorageItem {
				return []storage.StorageItem{
					{MetricName: "foo", Value: metric.Counter(246)},
				}
			},
		},
		{
			name: "update JSON metric: multiple values",
			req: testRequest{
				method: http.MethodPost,
				url:    "/updates/",
				body: bodyStringReader(`[
					{"id": "bar", "type": "gauge", "value": 123.123},
					{"id": "foo", "type": "counter", "delta": 1},
					{"id": "foo", "type": "counter", "delta": 1},
					{"id": "foo", "type": "counter", "delta": 1},
					{"id": "fuzz", "type": "counter", "delta": 2},
					{"id": "fuzz", "type": "counter", "delta": 2},
					{"id": "buzz", "type": "counter", "delta": 1}
				]`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `[]`,
				contentType: _http.ContentTypeApplicationJSON,
			},
			stgInitFunc: func(s storage.Storage) {
				s.SetCounterMetric("foo", 1)
			},
			stgCheckFunc: func() []storage.StorageItem {
				return []storage.StorageItem{
					{MetricName: "buzz", Value: metric.Counter(1)},
					{MetricName: "foo", Value: metric.Counter(4)},
					{MetricName: "fuzz", Value: metric.Counter(4)},
					{MetricName: "bar", Value: metric.Gauge(123.123)},
				}
			},
		},
		{
			name: "update JSON metric: multiple values, encrypted and gziped, with trusted subnet",
			req: testRequest{
				method: http.MethodPost,
				url:    "/updates/",
				body: bodyGzipReader(encryptString(`[
					{"id": "bar_enc", "type": "gauge", "value": 123.123},
					{"id": "foo_enc", "type": "counter", "delta": 1},
					{"id": "foo_enc", "type": "counter", "delta": 1},
					{"id": "foo_enc", "type": "counter", "delta": 1},
					{"id": "fuzz_enc", "type": "counter", "delta": 2},
					{"id": "fuzz_enc", "type": "counter", "delta": 2},
					{"id": "buzz_enc", "type": "counter", "delta": 1}
				]`)),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				encryption:  true,
				headers:     map[string][]string{"Accept-Encoding": {"gzip"}, "Content-Encoding": {"gzip"}, "X-Real-IP": {"192.168.0.1"}},
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `[]`,
				contentType: _http.ContentTypeApplicationJSON,
				headers:     map[string][]string{"Content-Encoding": {"gzip"}},
			},
			trustedSubnet: getIPNet("192.168.0.0/16"),
			stgInitFunc: func(s storage.Storage) {
				s.SetCounterMetric("foo_enc", 1)
			},
			stgCheckFunc: func() []storage.StorageItem {
				return []storage.StorageItem{
					{MetricName: "buzz_enc", Value: metric.Counter(1)},
					{MetricName: "foo_enc", Value: metric.Counter(4)},
					{MetricName: "fuzz_enc", Value: metric.Counter(4)},
					{MetricName: "bar_enc", Value: metric.Gauge(123.123)},
				}
			},
		},
	}

	runTests(t, tests)
}

func TestUpdateMetricJSONBatchInDb(t *testing.T) {
	tests := []testItem{
		{
			name: "update JSON metric: gauge",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/updates/",
				body:        bodyStringReader(`[{"id": "foo", "type": "gauge", "value": 123.0}]`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `[]`,
				contentType: _http.ContentTypeApplicationJSON,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().
					SetMetrics([]storage.StorageItem{{MetricName: "foo", Value: metric.Gauge(123.0)}}).
					Return(nil)
			},
		},
		{
			name: "update JSON metric: gauge db err",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/updates/",
				body:        bodyStringReader(`[{"id": "foo", "type": "gauge", "value": 123.0}]`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusInternalServerError,
				body:        http.StatusText(http.StatusInternalServerError),
				contentType: _http.ContentTypeTextPlain,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().
					SetMetrics([]storage.StorageItem{{MetricName: "foo", Value: metric.Gauge(123.0)}}).
					Return(errors.New("db error"))
			},
		},
		{
			name: "update JSON metric: counter",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/updates/",
				body:        bodyStringReader(`[{"id": "foo", "type": "counter", "delta": 123}]`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `[]`,
				contentType: _http.ContentTypeApplicationJSON,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().
					SetMetrics([]storage.StorageItem{{MetricName: "foo", Value: metric.Counter(123)}}).
					Return(nil)
			},
		},
		{
			name: "update JSON metric: counter db err",
			req: testRequest{
				method:      http.MethodPost,
				url:         "/updates/",
				body:        bodyStringReader(`[{"id": "foo", "type": "counter", "delta": 123}]`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
			},
			resp: testResponse{
				code:        http.StatusInternalServerError,
				body:        http.StatusText(http.StatusInternalServerError),
				contentType: _http.ContentTypeTextPlain,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().
					SetMetrics([]storage.StorageItem{{MetricName: "foo", Value: metric.Counter(123)}}).
					Return(errors.New("db error"))
			},
		},
	}

	runTests(t, tests)
}

func TestUpdateMetricJSONBatchWithHash(t *testing.T) {
	tests := []testItem{
		{
			name: "update JSON metric: multiple with correct hash",
			req: testRequest{
				method: http.MethodPost,
				url:    "/updates/",
				body: bodyStringReader(`[
					{"id":"Sys","type":"gauge","value":13220880,"hash":"48a93e5dde0297029bf66cc10a1cdda9be6f858667ea885dc1b0d810032aa292"},
					{"id":"PollCount","type":"counter","delta":5,"hash":"b9203cac5904e73da2504aabfb77a419d3d3f9a0baee3707c55070432c6ff5a8"},
					{"id":"PollCount","type":"counter","delta":5,"hash":"b9203cac5904e73da2504aabfb77a419d3d3f9a0baee3707c55070432c6ff5a8"}
				]`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				hmacKey:     strPointer("foobar"),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `[]`,
				contentType: _http.ContentTypeApplicationJSON,
			},
			stgCheckFunc: func() []storage.StorageItem {
				return []storage.StorageItem{
					{MetricName: "PollCount", Value: metric.Counter(10)},
					{MetricName: "Sys", Value: metric.Gauge(13220880)},
				}
			},
		},
		{
			name: "update JSON metric: multiple with incorrect hash",
			req: testRequest{
				method: http.MethodPost,
				url:    "/updates/",
				body: bodyStringReader(`[
					{"id":"Sys","type":"gauge","value":13220880,"hash":"48a93e5dde0297029bf66cc10a1cdda9be6f858667ea885dc1b0d810032aa292"},
					{"id":"PollCount","type":"counter","delta":5,"hash":"a9203cac5904e73da2504aabfb77a419d3d3f9a0baee3707c55070432c6ff5a8"},
					{"id":"PollCount","type":"counter","delta":5,"hash":"b9203cac5904e73da2504aabfb77a419d3d3f9a0baee3707c55070432c6ff5a8"}
				]`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				hmacKey:     strPointer("foobar"),
			},
			resp: testResponse{
				code:        http.StatusBadRequest,
				body:        http.StatusText(http.StatusBadRequest),
				contentType: _http.ContentTypeTextPlain,
			},
		},
	}
	runTests(t, tests)
}

func TestUpdateMetricJSONBatchWithHashInDb(t *testing.T) {
	tests := []testItem{
		{
			name: "update JSON metric: multiple with correct hash",
			req: testRequest{
				method: http.MethodPost,
				url:    "/updates/",
				body: bodyStringReader(`[
					{"id":"Sys","type":"gauge","value":13220880,"hash":"48a93e5dde0297029bf66cc10a1cdda9be6f858667ea885dc1b0d810032aa292"},
					{"id":"PollCount","type":"counter","delta":5,"hash":"b9203cac5904e73da2504aabfb77a419d3d3f9a0baee3707c55070432c6ff5a8"},
					{"id":"PollCount","type":"counter","delta":5,"hash":"b9203cac5904e73da2504aabfb77a419d3d3f9a0baee3707c55070432c6ff5a8"}
				]`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				hmacKey:     strPointer("foobar"),
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        `[]`,
				contentType: _http.ContentTypeApplicationJSON,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().
					SetMetrics([]storage.StorageItem{
						{MetricName: "Sys", Value: metric.Gauge(13220880)},
						{MetricName: "PollCount", Value: metric.Counter(5)},
						{MetricName: "PollCount", Value: metric.Counter(5)},
					}).
					Return(nil)
			},
		},
		{
			name: "update JSON metric: multiple with correct hash db error",
			req: testRequest{
				method: http.MethodPost,
				url:    "/updates/",
				body: bodyStringReader(`[
					{"id":"Sys","type":"gauge","value":13220880,"hash":"48a93e5dde0297029bf66cc10a1cdda9be6f858667ea885dc1b0d810032aa292"},
					{"id":"PollCount","type":"counter","delta":5,"hash":"b9203cac5904e73da2504aabfb77a419d3d3f9a0baee3707c55070432c6ff5a8"},
					{"id":"PollCount","type":"counter","delta":5,"hash":"b9203cac5904e73da2504aabfb77a419d3d3f9a0baee3707c55070432c6ff5a8"}
				]`),
				contentType: strPointer(_http.ContentTypeApplicationJSON),
				hmacKey:     strPointer("foobar"),
			},
			resp: testResponse{
				code:        http.StatusInternalServerError,
				body:        http.StatusText(http.StatusInternalServerError),
				contentType: _http.ContentTypeTextPlain,
			},
			dbStg: true,
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().
					SetMetrics([]storage.StorageItem{
						{MetricName: "Sys", Value: metric.Gauge(13220880)},
						{MetricName: "PollCount", Value: metric.Counter(5)},
						{MetricName: "PollCount", Value: metric.Counter(5)},
					}).
					Return(errors.New("db error"))
			},
		},
	}
	runTests(t, tests)
}

func ExampleMetricHandler_UpdateMetric_counter() {
	// Create HTTP client
	client := resty.New()

	// Create request object
	var v int64 = 123
	req := &metric.MetricsDTO{
		MType: "counter",
		ID:    "Counter1",
		Delta: &v,
	}
	var res metric.MetricsDTO

	// Send request
	_, err := client.R().
		SetBody(req).
		SetResult(&res).
		Post("http://localhost:8080/update/")

	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}

func ExampleMetricHandler_UpdateMetric_gauge() {
	// Create HTTP client
	client := resty.New()

	// Create request object
	v := 123.123
	req := &metric.MetricsDTO{
		MType: "gauge",
		ID:    "Gauge1",
		Value: &v,
	}
	var res metric.MetricsDTO

	// Send request
	_, err := client.R().
		SetBody(req).
		SetResult(&res).
		Post("http://localhost:8080/update/")

	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
