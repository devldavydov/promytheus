package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateMetricsHandler(t *testing.T) {
	type response struct {
		code        int
		body        string
		contentType string
	}

	type request struct {
		method string
		url    string
	}

	tests := []struct {
		name string
		req  request
		resp response
	}{
		{
			name: "Failed GET request",
			req: request{
				method: http.MethodGet,
				url:    "/update/gauge/metric/1",
			},
			resp: response{
				code:        http.StatusMethodNotAllowed,
				body:        "",
				contentType: "",
			},
		},
		{
			name: "Incorrect URL: parts count #1",
			req: request{
				method: http.MethodPost,
				url:    "/update/gauge/metric1/1/2/3",
			},
			resp: response{
				code:        http.StatusNotFound,
				body:        "404 page not found\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Incorrect URL: parts count #2",
			req: request{
				method: http.MethodPost,
				url:    "/update/gauge/metric1",
			},
			resp: response{
				code:        http.StatusNotFound,
				body:        "404 page not found\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Incorrect URL: unknown metric type",
			req: request{
				method: http.MethodPost,
				url:    "/update/fuzz/metric1/1",
			},
			resp: response{
				code:        http.StatusNotImplemented,
				body:        "Not Implemented",
				contentType: "text/plain",
			},
		},
		{
			name: "Incorrect URL: empty metric name",
			req: request{
				method: http.MethodPost,
				url:    "/update/gauge//1",
			},
			resp: response{
				code:        http.StatusBadRequest,
				body:        "Bad Request",
				contentType: "text/plain",
			},
		},
		{
			name: "Incorrect URL: incorrect gauge val, #1",
			req: request{
				method: http.MethodPost,
				url:    "/update/gauge/metric1/foobar",
			},
			resp: response{
				code:        http.StatusBadRequest,
				body:        "Bad Request",
				contentType: "text/plain",
			},
		},
		{
			name: "Incorrect URL: incorrect counter val, #1",
			req: request{
				method: http.MethodPost,
				url:    "/update/counter/metric1/foobar",
			},
			resp: response{
				code:        http.StatusBadRequest,
				body:        "Bad Request",
				contentType: "text/plain",
			},
		},
		{
			name: "Incorrect URL: incorrect counter val, #2",
			req: request{
				method: http.MethodPost,
				url:    "/update/counter/metric1/1.234",
			},
			resp: response{
				code:        http.StatusBadRequest,
				body:        "Bad Request",
				contentType: "text/plain",
			},
		},
		{
			name: "Incorrect URL: incorrect counter val, #3",
			req: request{
				method: http.MethodPost,
				url:    "/update/counter/metric1/-1234",
			},
			resp: response{
				code:        http.StatusBadRequest,
				body:        "Bad Request",
				contentType: "text/plain",
			},
		},
		{
			name: "Correct gauge",
			req: request{
				method: http.MethodPost,
				url:    "/update/gauge/metric1/1.234",
			},
			resp: response{
				code:        http.StatusOK,
				body:        "OK",
				contentType: "text/plain",
			},
		},
		{
			name: "Correct counter",
			req: request{
				method: http.MethodPost,
				url:    "/update/counter/metric1/1234",
			},
			resp: response{
				code:        http.StatusOK,
				body:        "OK",
				contentType: "text/plain",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricsHandler := NewMetricsHandler(storage.NewMemStorage(), logrus.New())
			r := NewRouter(metricsHandler)
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, contentType, body := testRequest(t, ts, tt.req.method, tt.req.url)

			assert.Equal(t, tt.resp.code, statusCode)
			assert.Equal(t, tt.resp.body, body)
			assert.Equal(t, tt.resp.contentType, contentType)
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (int, string, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp.StatusCode, resp.Header.Get("Content-Type"), string(respBody)
}
