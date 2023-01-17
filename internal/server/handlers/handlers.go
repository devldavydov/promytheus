package handlers

import "net/http"

const (
	UpdateMetricsUrlPattern string = "/update/"
)

type Handler interface {
	Handle(handleFunc func(pattern string, handler func(http.ResponseWriter, *http.Request)))
}
