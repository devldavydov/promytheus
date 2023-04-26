package metric

import (
	"net/http"

	_http "github.com/devldavydov/promytheus/internal/common/http"
)

// Ping checks storage connection.
//
//	@Summary	Check storage connection
//	@Produce	plain/text
//	@Success	200	"Check successful"
//	@Failure	500	"Internal error"
//	@Router		/ping [get]
func (handler *MetricHandler) Ping(rw http.ResponseWriter, req *http.Request) {
	if handler.storage.Ping() {
		_http.CreateStatusResponse(rw, http.StatusOK)
		return
	}
	_http.CreateStatusResponse(rw, http.StatusInternalServerError)
}
