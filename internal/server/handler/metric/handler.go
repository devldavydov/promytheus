// Package metric is a package for metrics HTTP handlers.
package metric

import (
	"errors"
	"net"
	"net/http"

	_http "github.com/devldavydov/promytheus/internal/common/http"
	_middleware "github.com/devldavydov/promytheus/internal/server/middleware"
	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

/*
@Title Metrics API
@Version 1.0
@Description Yandex Praktikum metrics server API

@Host      localhost:8080
@BasePath  /
*/

var (
	ErrUnknownMetricType = errors.New("unknowm metric type")
	ErrEmptyMetricName   = errors.New("empty metric name")
	ErrWrongMetricValue  = errors.New("wrong metric value")
	ErrMetricHashCheck   = errors.New("metric hash check fail")
)

type MetricHandler struct {
	storage storage.Storage
	hmacKey *string
	logger  *logrus.Logger
}

func NewHandler(
	router chi.Router,
	storage storage.Storage,
	hmacKey *string,
	trustedSubnet *net.IPNet,
	logger *logrus.Logger,
) *MetricHandler {
	handler := &MetricHandler{storage: storage, hmacKey: hmacKey, logger: logger}

	mdlwrTrusted := _middleware.NewTrusted(trustedSubnet)

	router.Group(func(r chi.Router) {
		r.Use(mdlwrTrusted.Handle)

		r.Post("/update/{metricType}/{metricName}/{metricValue}", handler.UpdateMetric)
		r.Post("/update/", handler.UpdateMetricJSON)
		r.Post("/updates/", handler.UpdateMetricJSONBatch)
	})

	router.Get("/value/{metricType}/{metricName}", handler.GetMetric)
	router.Post("/value/", handler.GetMetricJSON)
	router.Get("/", handler.GetMetrics)
	router.Get("/ping", handler.Ping)

	return handler
}

func CreateResponseOnRequestError(rw http.ResponseWriter, err error) {
	if errors.Is(err, ErrUnknownMetricType) {
		_http.CreateStatusResponse(rw, http.StatusNotImplemented)
		return
	}
	if errors.Is(err, ErrMetricHashCheck) {
		_http.CreateStatusResponse(rw, http.StatusBadRequest)
		return
	}

	_http.CreateStatusResponse(rw, http.StatusBadRequest)
}
