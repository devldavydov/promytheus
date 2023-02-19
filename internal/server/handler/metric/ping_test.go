package metric

import (
	"net/http"
	"testing"

	_http "github.com/devldavydov/promytheus/internal/common/http"
	"github.com/devldavydov/promytheus/internal/server/mocks"
)

func TestPing(t *testing.T) {
	tests := []testItem{
		{
			name: "ping db connection: failed POST request",
			req: testRequest{
				method: http.MethodPost,
				url:    "/ping",
			},
			resp: testResponse{
				code:        http.StatusMethodNotAllowed,
				body:        "",
				contentType: "",
			},
		},
		{
			name: "ping db connection: success check",
			req: testRequest{
				method: http.MethodGet,
				url:    "/ping",
			},
			resp: testResponse{
				code:        http.StatusOK,
				body:        http.StatusText(http.StatusOK),
				contentType: _http.ContentTypeTextPlain,
			},
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().Ping().Return(true)
			},
		},
		{
			name: "ping db connection: fail check",
			req: testRequest{
				method: http.MethodGet,
				url:    "/ping",
			},
			resp: testResponse{
				code:        http.StatusInternalServerError,
				body:        http.StatusText(http.StatusInternalServerError),
				contentType: _http.ContentTypeTextPlain,
			},
			stgMockFunc: func(ms *mocks.MockStorage) {
				ms.EXPECT().Ping().Return(false)
			},
		},
	}

	runTests(t, tests)
}
