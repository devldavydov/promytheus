package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/devldavydov/promytheus/internal/common/types"
	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/sirupsen/logrus"
)

type UpdateMetricsHandler struct {
	urlPattern string
	storage    storage.Storage
	logger     *logrus.Logger
}

type requestParams struct {
	metricType   string
	metricName   string
	gaugeValue   types.Gauge
	counterValue types.Counter
}

func NewUpdateMetricsHandler(urlPattern string, storage storage.Storage, logger *logrus.Logger) *UpdateMetricsHandler {
	return &UpdateMetricsHandler{urlPattern: urlPattern, storage: storage, logger: logger}
}

func (handler *UpdateMetricsHandler) Handle(handleFunc func(pattern string, handler func(http.ResponseWriter, *http.Request))) {
	handleFunc(handler.urlPattern, handler.HandlerFunc())
}

func (handler *UpdateMetricsHandler) HandlerFunc() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			handler.createResponse(rw, "text/plain", http.StatusMethodNotAllowed, "Method Not Allowed")
			return
		}

		params, err := handler.parseRequest(req)
		if err != nil {
			handler.logger.Errorf("Incorrect update request [%s], err: %v", req.URL, err)
			handler.createResponse(rw, "text/plain", http.StatusBadRequest, "Bad Request")
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

func (handler *UpdateMetricsHandler) parseRequest(req *http.Request) (requestParams, error) {
	url := strings.TrimPrefix(req.URL.Path, handler.urlPattern)
	parts := strings.Split(url, "/")

	if len(parts) != 3 {
		return requestParams{}, errors.New("incorrect url parts count")
	}

	if !types.AllTypes[parts[0]] {
		return requestParams{}, errors.New("unknown metric type")
	}

	if len(parts[1]) == 0 {
		return requestParams{}, errors.New("empty metric name")
	}

	if types.GaugeTypeName == parts[0] {
		gaugeVal, err := types.NewGaugeFromString(parts[2])
		if err != nil {
			return requestParams{}, fmt.Errorf("incorrect %s val", types.GaugeTypeName)
		}
		return requestParams{
			metricType: types.GaugeTypeName,
			metricName: parts[1],
			gaugeValue: gaugeVal,
		}, nil
	} else if types.CounterTypeName == parts[0] {
		counterVal, err := types.NewCounterFromString(parts[2])
		if err != nil {
			return requestParams{}, fmt.Errorf("incorrect %s val", types.CounterTypeName)
		}
		return requestParams{
			metricType:   types.CounterTypeName,
			metricName:   parts[1],
			counterValue: counterVal,
		}, nil
	}

	return requestParams{}, nil
}

func (handler *UpdateMetricsHandler) createResponse(rw http.ResponseWriter, contentType string, statusCode int, body string) {
	rw.Header().Set("Content-Type", contentType)
	rw.WriteHeader(statusCode)
	io.WriteString(rw, body)
}
