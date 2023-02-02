package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var (
	ErrUnknownMetricType = errors.New("unknowm metric type")
	ErrEmptyMetricName   = errors.New("empty metric name")
	ErrWrongMetricValue  = errors.New("wrong metric value")
)

const (
	ContentTypeApplicationJSON string = "application/json; charset=utf-8"
	ContentTypeTextPlain       string = "text/plain; charset=utf-8"
	ContentTypeHTML            string = "text/html; charset=utf-8"
)

func createResponse(rw http.ResponseWriter, contentType string, statusCode int, body string) {
	rw.Header().Set("Content-Type", contentType)
	rw.WriteHeader(statusCode)
	io.WriteString(rw, body)
}

func createJSONResponse(rw http.ResponseWriter, statusCode int, body interface{}) {
	rw.Header().Set("Content-Type", ContentTypeApplicationJSON)
	rw.WriteHeader(statusCode)
	json.NewEncoder(rw).Encode(body)
}
