package handlers

import "net/http"

const (
	UpdateMetricsURLPattern string = "/update/"
)

type Handler interface {
	Handle(handleFunc func(pattern string, handler func(http.ResponseWriter, *http.Request)))
}
