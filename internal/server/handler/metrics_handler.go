package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"text/template"

	_http "github.com/devldavydov/promytheus/internal/common/http"
	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type MetricsHandler struct {
	storage storage.Storage
	hmacKey *string
	logger  *logrus.Logger
}

func NewMetricsHandler(storage storage.Storage, hmacKey *string, logger *logrus.Logger) *MetricsHandler {
	return &MetricsHandler{storage: storage, hmacKey: hmacKey, logger: logger}
}

func (handler *MetricsHandler) UpdateMetric(rw http.ResponseWriter, req *http.Request) {
	params, err := handler.parseUpdateRequest(chi.URLParam(req, "metricType"), chi.URLParam(req, "metricName"), chi.URLParam(req, "metricValue"))
	if err != nil {
		handler.logger.Errorf("Incorrect update metric request [%s], err: %v", req.URL, err)
		createResponseOnRequestError(rw, err)
		return
	}

	if metric.GaugeTypeName == params.metricType {
		err = handler.storage.SetGaugeMetric(params.metricName, params.gaugeValue)
	} else if metric.CounterTypeName == params.metricType {
		err = handler.storage.SetCounterMetric(params.metricName, params.counterValue)
	}

	if err != nil {
		handler.logger.Errorf("Update metric error on request [%s], err: %v", req.URL, err)
		createErrResponse(rw, http.StatusInternalServerError)
		return
	}

	createResponse(rw, _http.ContentTypeTextPlain, http.StatusOK, http.StatusText(http.StatusOK))
}

func (handler *MetricsHandler) UpdateMetricJSON(rw http.ResponseWriter, req *http.Request) {
	var metricReq metric.MetricsDTO

	err := json.NewDecoder(req.Body).Decode(&metricReq)
	if err != nil {
		createErrResponse(rw, http.StatusBadRequest)
		return
	}

	params, err := handler.parseUpdateRequestJSON(metricReq)
	if err != nil {
		handler.logger.Errorf("Incorrect update metric request [%s], JSON: [%v] , err: %v", req.URL, metricReq, err)
		createResponseOnRequestError(rw, err)
		return
	}

	metricResp := metric.MetricsDTO{ID: metricReq.ID, MType: metricReq.MType}
	var val interface{}

	if metric.GaugeTypeName == params.metricType {
		val, err = handler.storage.SetAndGetGaugeMetric(params.metricName, params.gaugeValue)
	} else if metric.CounterTypeName == params.metricType {
		val, err = handler.storage.SetAndGetCounterMetric(params.metricName, params.counterValue)
	}

	if err != nil {
		handler.logger.Errorf("Update metric error on request [%s], JSON: [%v], err: %v", req.URL, metricReq, err)
		createErrResponse(rw, http.StatusInternalServerError)
		return
	}

	if metric.GaugeTypeName == params.metricType {
		metricResp.Value = val.(metric.Gauge).FloatP()
	} else if metric.CounterTypeName == params.metricType {
		metricResp.Delta = val.(metric.Counter).IntP()
	}

	if handler.hmacKey != nil {
		hash := val.(metric.MetricValue).Hmac(params.metricName, *handler.hmacKey)
		metricResp.Hash = &hash
	}

	createJSONResponse(rw, http.StatusOK, metricResp)
}

func (handler *MetricsHandler) GetMetric(rw http.ResponseWriter, req *http.Request) {
	metricType, metricName := chi.URLParam(req, "metricType"), chi.URLParam(req, "metricName")

	err := handler.checkMetricsCommon(metricType, metricName)
	if err != nil {
		handler.logger.Errorf("Incorrect get metric request [%s], err: %v", req.URL, err)
		createResponseOnRequestError(rw, err)
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
			createErrResponse(rw, http.StatusNotFound)
			return
		}

		handler.logger.Errorf("Get metric error on request [%s], err: %v", req.URL, err)
		createErrResponse(rw, http.StatusInternalServerError)
		return
	}

	createResponse(rw, _http.ContentTypeTextPlain, http.StatusOK, value.String())
}

func (handler *MetricsHandler) GetMetricJSON(rw http.ResponseWriter, req *http.Request) {
	var metricReq metric.MetricsDTO

	err := json.NewDecoder(req.Body).Decode(&metricReq)
	if err != nil {
		createErrResponse(rw, http.StatusBadRequest)
		return
	}

	err = handler.checkMetricsCommon(metricReq.MType, metricReq.ID)
	if err != nil {
		handler.logger.Errorf("Incorrect get metric request [%s], JSON: [%v], err: %v", req.URL, metricReq, err)
		createResponseOnRequestError(rw, err)
		return
	}

	metricResp := metric.MetricsDTO{ID: metricReq.ID, MType: metricReq.MType}
	var val interface{}

	if metric.GaugeTypeName == metricReq.MType {
		val, err = handler.storage.GetGaugeMetric(metricReq.ID)
	} else if metric.CounterTypeName == metricReq.MType {
		val, err = handler.storage.GetCounterMetric(metricReq.ID)
	}

	if err != nil {
		if errors.Is(err, storage.ErrMetricNotFound) {
			createErrResponse(rw, http.StatusNotFound)
			return
		}

		handler.logger.Errorf("Get metric error on request [%s], JSON: [%v] err: %v", req.URL, metricReq, err)
		createErrResponse(rw, http.StatusInternalServerError)
		return
	}

	if metric.GaugeTypeName == metricReq.MType {
		metricResp.Value = val.(metric.Gauge).FloatP()
	} else if metric.CounterTypeName == metricReq.MType {
		metricResp.Delta = val.(metric.Counter).IntP()
	}

	if handler.hmacKey != nil {
		hash := val.(metric.MetricValue).Hmac(metricReq.ID, *handler.hmacKey)
		metricResp.Hash = &hash
	}

	createJSONResponse(rw, http.StatusOK, metricResp)
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
		createErrResponse(rw, http.StatusInternalServerError)
		return
	}
	tmpl, _ := template.New("metrics").Parse(pageTemplate)
	buf := new(bytes.Buffer)
	tmpl.Execute(buf, metrics)
	createResponse(rw, _http.ContentTypeHTML, http.StatusOK, buf.String())
}

func createResponseOnRequestError(rw http.ResponseWriter, err error) {
	if errors.Is(err, ErrUnknownMetricType) {
		createErrResponse(rw, http.StatusNotImplemented)
		return
	}
	if errors.Is(err, ErrMetricHashCheck) {
		createErrResponse(rw, http.StatusBadRequest)
		return
	}

	createErrResponse(rw, http.StatusBadRequest)
}
