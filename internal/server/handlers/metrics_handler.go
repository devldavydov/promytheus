package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/devldavydov/promytheus/internal/common/types"
	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type MetricsHandler struct {
	storage storage.Storage
	logger  *logrus.Logger
}

type requestParams struct {
	metricType   string
	metricName   string
	gaugeValue   types.Gauge
	counterValue types.Counter
}

func NewMetricsHandler(storage storage.Storage, logger *logrus.Logger) *MetricsHandler {
	return &MetricsHandler{storage: storage, logger: logger}
}

func (handler *MetricsHandler) UpdateMetrics() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		params, err := handler.parseRequest(chi.URLParam(req, "metricType"), chi.URLParam(req, "metricName"), chi.URLParam(req, "metricValue"))
		if err != nil {
			handler.logger.Errorf("Incorrect update request [%s], err: %v", req.URL, err)

			if errors.As(err, &IncorrectURLUnknownMetricTypeP) {
				handler.createResponse(rw, "text/plain", http.StatusNotImplemented, "Not Implemented")
			} else {
				handler.createResponse(rw, "text/plain", http.StatusBadRequest, "Bad Request")
			}
			return
		}

		if types.GaugeTypeName == params.metricType {
			handler.storage.SetGaugeMetric(params.metricName, params.gaugeValue)
		} else if types.CounterTypeName == params.metricType {
			handler.storage.SetCounterMetric(params.metricName, params.counterValue)
		}
		handler.createResponse(rw, "text/plain", http.StatusOK, "OK")
	}
}

func (handler *MetricsHandler) parseRequest(metricType, metricName, metricValue string) (requestParams, error) {
	if !types.AllTypes[metricType] {
		return requestParams{}, &IncorrectURLUnknownMetricType{err: "unknown metric type"}
	}

	if len(metricName) == 0 {
		return requestParams{}, &IncorrectURLEmptyMetricName{err: "empty metric name"}
	}

	if types.GaugeTypeName == metricType {
		gaugeVal, err := types.NewGaugeFromString(metricValue)
		if err != nil {
			return requestParams{}, &IncorrectURLWrongMetricValue{err: fmt.Sprintf("incorrect %s val", types.GaugeTypeName)}
		}
		return requestParams{
			metricType: types.GaugeTypeName,
			metricName: metricName,
			gaugeValue: gaugeVal,
		}, nil
	} else if types.CounterTypeName == metricType {
		counterVal, err := types.NewCounterFromString(metricValue)
		if err != nil {
			return requestParams{}, &IncorrectURLWrongMetricValue{fmt.Sprintf("incorrect %s val", types.CounterTypeName)}
		}
		return requestParams{
			metricType:   types.CounterTypeName,
			metricName:   metricName,
			counterValue: counterVal,
		}, nil
	}

	return requestParams{}, nil
}

func (handler *MetricsHandler) createResponse(rw http.ResponseWriter, contentType string, statusCode int, body string) {
	rw.Header().Set("Content-Type", contentType)
	rw.WriteHeader(statusCode)
	io.WriteString(rw, body)
}
