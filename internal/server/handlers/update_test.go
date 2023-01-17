package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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
				body:        "Method Not Allowed",
				contentType: "text/plain",
			},
		},
		{
			name: "Incorrect URL: parts count",
			req: request{
				method: http.MethodPost,
				url:    "/update/gauge/metric1/1/2/3",
			},
			resp: response{
				code:        http.StatusBadRequest,
				body:        "Bad Request",
				contentType: "text/plain",
			},
		},
		{
			name: "Incorrect URL: unknown metric type",
			req: request{
				method: http.MethodPost,
				url:    "/update/fuzz/metric1/1",
			},
			resp: response{
				code:        http.StatusBadRequest,
				body:        "Bad Request",
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
			handler := NewUpdateMetricsHandler("/update/", storage.NewMemStorage(), logrus.New())

			request := httptest.NewRequest(tt.req.method, tt.req.url, nil)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.HandlerFunc())
			h.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.resp.code, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.resp.body, string(resBody))
			assert.Equal(t, tt.resp.contentType, res.Header.Get("Content-Type"))
		})
	}
}
