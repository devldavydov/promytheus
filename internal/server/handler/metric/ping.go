package metric

import (
	"net/http"

	_http "github.com/devldavydov/promytheus/internal/common/http"
)

func (handler *MetricHandler) Ping(rw http.ResponseWriter, req *http.Request) {
	if handler.dbstorage.Ping() {
		_http.CreateStatusResponse(rw, http.StatusOK)
		return
	}
	_http.CreateStatusResponse(rw, http.StatusInternalServerError)
}
