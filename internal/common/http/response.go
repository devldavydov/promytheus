package http

import (
	"encoding/json"
	"io"
	"net/http"
)

func CreateResponse(rw http.ResponseWriter, contentType string, statusCode int, body string) {
	rw.Header().Set("Content-Type", contentType)
	rw.WriteHeader(statusCode)
	io.WriteString(rw, body)
}

func CreateStatusResponse(rw http.ResponseWriter, statusCode int) {
	rw.Header().Set("Content-Type", ContentTypeTextPlain)
	rw.WriteHeader(statusCode)
	io.WriteString(rw, http.StatusText(statusCode))
}

func CreateJSONResponse(rw http.ResponseWriter, statusCode int, body interface{}) {
	rw.Header().Set("Content-Type", ContentTypeApplicationJSON)
	rw.WriteHeader(statusCode)
	json.NewEncoder(rw).Encode(body)
}
