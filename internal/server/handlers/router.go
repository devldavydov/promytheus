package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(metricsHanlder *MetricsHandler, middlewares ...func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()
	r.Use(middlewares...)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", metricsHanlder.UpdateMetric())
	r.Get("/value/{metricType}/{metricName}", metricsHanlder.GetMetric())
	r.Get("/", metricsHanlder.GetMetrics())
	return r
}
