package metric

import (
	"net/http"
	"testing"

	_http "github.com/devldavydov/promytheus/internal/common/http"
	"github.com/devldavydov/promytheus/internal/server/storage"
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
				body:        `{"id":"foo","type":"gauge","value":123}` + "\n",
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
				body:        `{"id":"foo","type":"counter","delta":123}` + "\n",
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
				body:        `{"id":"foo","type":"counter","delta":246}` + "\n",
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
				body:        `{"id":"foo","type":"counter","delta":246}` + "\n",
				contentType: _http.ContentTypeApplicationJSON,
				headers:     map[string][]string{"Content-Encoding": {"gzip"}},
			},
			stgInitFunc: func(s storage.Storage) {
				s.SetCounterMetric("foo", 123)
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
				body:        `{"id":"Sys","type":"gauge","value":13220880,"hash":"48a93e5dde0297029bf66cc10a1cdda9be6f858667ea885dc1b0d810032aa292"}` + "\n",
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
				body:        `{"id":"PollCount","type":"counter","delta":5,"hash":"b9203cac5904e73da2504aabfb77a419d3d3f9a0baee3707c55070432c6ff5a8"}` + "\n",
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
