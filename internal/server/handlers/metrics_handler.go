package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"text/template"

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

func (handler *MetricsHandler) UpdateMetric() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		params, err := handler.parseUpdateRequest(chi.URLParam(req, "metricType"), chi.URLParam(req, "metricName"), chi.URLParam(req, "metricValue"))
		if err != nil {
			handler.logger.Errorf("Incorrect update request [%s], err: %v", req.URL, err)

			if errors.As(err, &IncorrectURLUnknownMetricTypeP) {
				handler.createResponse(rw, ContentTypeTextPlain, http.StatusNotImplemented, ResponseNotImplemented)
			} else {
				handler.createResponse(rw, ContentTypeTextPlain, http.StatusBadRequest, ResponseBadRequest)
			}
			return
		}

		if types.GaugeTypeName == params.metricType {
			handler.storage.SetGaugeMetric(params.metricName, params.gaugeValue)
		} else if types.CounterTypeName == params.metricType {
			handler.storage.SetCounterMetric(params.metricName, params.counterValue)
		}
		handler.createResponse(rw, ContentTypeTextPlain, http.StatusOK, ResponseOk)
	}
}

func (handler *MetricsHandler) GetMetric() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		metricType, metricName := chi.URLParam(req, "metricType"), chi.URLParam(req, "metricName")

		err := handler.checkGetRequest(metricType, metricName)
		if err != nil {
			handler.logger.Errorf("Incorrect get request [%s], err: %v", req.URL, err)

			if errors.As(err, &IncorrectURLUnknownMetricTypeP) {
				handler.createResponse(rw, ContentTypeTextPlain, http.StatusNotImplemented, ResponseNotImplemented)
			} else {
				handler.createResponse(rw, ContentTypeTextPlain, http.StatusBadRequest, ResponseBadRequest)
			}
			return
		}

		var value fmt.Stringer
		if types.GaugeTypeName == metricType {
			value, err = handler.storage.GetGaugeMetric(metricName)
		} else if types.CounterTypeName == metricType {
			value, err = handler.storage.GetCounterMetric(metricName)
		}

		if err != nil {
			if errors.As(err, &storage.MetricNotFoundErrorP) {
				handler.createResponse(rw, ContentTypeTextPlain, http.StatusNotFound, ResponseNotFound)
			} else {
				handler.logger.Errorf("Get metric error on request [%s], err: %v", req.URL, err)
				handler.createResponse(rw, ContentTypeTextPlain, http.StatusInternalServerError, ResponseInternalServerError)
			}
		} else {
			handler.createResponse(rw, ContentTypeTextPlain, http.StatusOK, value.String())
		}
	}
}

func (handler *MetricsHandler) GetMetrics() http.HandlerFunc {
	pageTemplate := `
	<html>
		<body>
			<table border="1">
				<tr>
					<th>Metric Type</th>
					<th>Metric Name</th>
					<th>Metric Value</th<
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

	return func(rw http.ResponseWriter, req *http.Request) {
		metrics, err := handler.storage.GetAllMetrics()
		if err != nil {
			handler.logger.Errorf("Get all metrics error: %v", err)
			handler.createResponse(rw, ContentTypeTextPlain, http.StatusInternalServerError, ResponseInternalServerError)
			return
		}
		tmpl, _ := template.New("metrics").Parse(pageTemplate)
		handler.createResponseTmpl(rw, ContentTypeHtml, http.StatusOK, func(w io.Writer) {
			tmpl.Execute(w, metrics)
		})
	}
}

func (handler *MetricsHandler) parseUpdateRequest(metricType, metricName, metricValue string) (requestParams, error) {
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

func (handler *MetricsHandler) checkGetRequest(metricType, metricName string) error {
	if !types.AllTypes[metricType] {
		return &IncorrectURLUnknownMetricType{err: "unknown metric type"}
	}
	return nil
}

func (handler *MetricsHandler) createResponse(rw http.ResponseWriter, contentType string, statusCode int, body string) {
	rw.Header().Set("Content-Type", contentType)
	rw.WriteHeader(statusCode)
	io.WriteString(rw, body)
}

func (handler *MetricsHandler) createResponseTmpl(rw http.ResponseWriter, contentType string, statusCode int, tmpl func(io.Writer)) {
	rw.Header().Set("Content-Type", contentType)
	rw.WriteHeader(statusCode)
	tmpl(rw)
}
