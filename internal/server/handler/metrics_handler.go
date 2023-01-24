package handler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/devldavydov/promytheus/internal/common/metric"
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
	gaugeValue   metric.Gauge
	counterValue metric.Counter
}

func NewMetricsHandler(storage storage.Storage, logger *logrus.Logger) *MetricsHandler {
	return &MetricsHandler{storage: storage, logger: logger}
}

func (handler *MetricsHandler) UpdateMetric(rw http.ResponseWriter, req *http.Request) {
	params, err := handler.parseUpdateRequest(chi.URLParam(req, "metricType"), chi.URLParam(req, "metricName"), chi.URLParam(req, "metricValue"))
	if err != nil {
		handler.logger.Errorf("Incorrect update request [%s], err: %v", req.URL, err)

		if errors.Is(err, ErrUnknownMetricType) {
			handler.createResponse(rw, ContentTypeTextPlain, http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
		} else {
			handler.createResponse(rw, ContentTypeTextPlain, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		}
		return
	}

	if metric.GaugeTypeName == params.metricType {
		handler.storage.SetGaugeMetric(params.metricName, params.gaugeValue)
	} else if metric.CounterTypeName == params.metricType {
		handler.storage.SetCounterMetric(params.metricName, params.counterValue)
	}
	handler.createResponse(rw, ContentTypeTextPlain, http.StatusOK, http.StatusText(http.StatusOK))
}

func (handler *MetricsHandler) GetMetric(rw http.ResponseWriter, req *http.Request) {
	metricType, metricName := chi.URLParam(req, "metricType"), chi.URLParam(req, "metricName")

	err := handler.checkGetRequest(metricType)
	if err != nil {
		handler.logger.Errorf("Incorrect get request [%s], err: %v", req.URL, err)

		if errors.Is(err, ErrUnknownMetricType) {
			handler.createResponse(rw, ContentTypeTextPlain, http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
		} else {
			handler.createResponse(rw, ContentTypeTextPlain, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		}
		return
	}

	var value fmt.Stringer
	if metric.GaugeTypeName == metricType {
		value, err = handler.storage.GetGaugeMetric(metricName)
	} else if metric.CounterTypeName == metricType {
		value, err = handler.storage.GetCounterMetric(metricName)
	}

	if err != nil {
		if errors.Is(err, storage.ErrMetricNotFound) {
			handler.createResponse(rw, ContentTypeTextPlain, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		} else {
			handler.logger.Errorf("Get metric error on request [%s], err: %v", req.URL, err)
			handler.createResponse(rw, ContentTypeTextPlain, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
	} else {
		handler.createResponse(rw, ContentTypeTextPlain, http.StatusOK, value.String())
	}
}

func (handler *MetricsHandler) GetMetrics(rw http.ResponseWriter, req *http.Request) {
	pageTemplate := `
	<html>
		<body>
			<table border="1">
				<tr>
					<th>Metric Type</th>
					<th>Metric Name</th>
					<th>Metric Value</th>
				</tr>
				{{ range . }}
				<tr>
					<td>{{ .Value.TypeName }}</td>
					<td>{{ .MetricName }}</td>
					<td>{{ .Value }}</td>
				</tr>
				{{ end }}
			</table>
		</body>
	</html>
	`

	metrics, err := handler.storage.GetAllMetrics()
	if err != nil {
		handler.logger.Errorf("Get all metrics error: %v", err)
		handler.createResponse(rw, ContentTypeTextPlain, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	tmpl, _ := template.New("metrics").Parse(pageTemplate)
	buf := new(bytes.Buffer)
	tmpl.Execute(buf, metrics)
	handler.createResponse(rw, ContentTypeHTML, http.StatusOK, buf.String())
}

func (handler *MetricsHandler) parseUpdateRequest(metricType, metricName, metricValue string) (requestParams, error) {
	if !metric.AllTypes[metricType] {
		return requestParams{}, ErrUnknownMetricType
	}

	if len(metricName) == 0 {
		return requestParams{}, ErrEmptyMetricName
	}

	if metric.GaugeTypeName == metricType {
		gaugeVal, err := metric.NewGaugeFromString(metricValue)
		if err != nil {
			return requestParams{}, fmt.Errorf("incorrect %s: %w", metric.GaugeTypeName, ErrWrongMetricValue)
		}
		return requestParams{
			metricType: metric.GaugeTypeName,
			metricName: metricName,
			gaugeValue: gaugeVal,
		}, nil
	} else if metric.CounterTypeName == metricType {
		counterVal, err := metric.NewCounterFromString(metricValue)
		if err != nil {
			return requestParams{}, fmt.Errorf("incorrect %s: %w", metric.CounterTypeName, ErrWrongMetricValue)
		}
		return requestParams{
			metricType:   metric.CounterTypeName,
			metricName:   metricName,
			counterValue: counterVal,
		}, nil
	}

	return requestParams{}, nil
}

func (handler *MetricsHandler) checkGetRequest(metricType string) error {
	if !metric.AllTypes[metricType] {
		return ErrUnknownMetricType
	}
	return nil
}

func (handler *MetricsHandler) createResponse(rw http.ResponseWriter, contentType string, statusCode int, body string) {
	rw.Header().Set("Content-Type", contentType)
	rw.WriteHeader(statusCode)
	io.WriteString(rw, body)
}
