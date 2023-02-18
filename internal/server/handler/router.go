package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(metricsHanlder *MetricsHandler, middlewares ...func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()
	r.Use(middlewares...)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", metricsHanlder.UpdateMetric)
	r.Post("/update/", metricsHanlder.UpdateMetricJSON)
	r.Get("/value/{metricType}/{metricName}", metricsHanlder.GetMetric)
	r.Post("/value/", metricsHanlder.GetMetricJSON)
	r.Get("/", metricsHanlder.GetMetrics)
	r.Get("/ping", metricsHanlder.Ping)
	return r
}
