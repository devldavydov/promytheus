package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	_http "github.com/devldavydov/promytheus/internal/common/http"
)

var (
	ErrUnknownMetricType = errors.New("unknowm metric type")
	ErrEmptyMetricName   = errors.New("empty metric name")
	ErrWrongMetricValue  = errors.New("wrong metric value")
	ErrMetricHashCheck   = errors.New("metric hash check fail")
)

func createResponse(rw http.ResponseWriter, contentType string, statusCode int, body string) {
	rw.Header().Set("Content-Type", contentType)
	rw.WriteHeader(statusCode)
	io.WriteString(rw, body)
}

func createJSONResponse(rw http.ResponseWriter, statusCode int, body interface{}) {
	rw.Header().Set("Content-Type", _http.ContentTypeApplicationJSON)
	rw.WriteHeader(statusCode)
	json.NewEncoder(rw).Encode(body)
}
