package metric

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	_http "github.com/devldavydov/promytheus/internal/common/http"
	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

func (handler *MetricHandler) GetMetric(rw http.ResponseWriter, req *http.Request) {
	metricType, metricName := chi.URLParam(req, "metricType"), chi.URLParam(req, "metricName")

	err := handler.checkMetricsCommon(metricType, metricName)
	if err != nil {
		handler.logger.Errorf("Incorrect get metric request [%s], err: %v", req.URL, err)
		CreateResponseOnRequestError(rw, err)
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
			_http.CreateStatusResponse(rw, http.StatusNotFound)
			return
		}

		handler.logger.Errorf("Get metric error on request [%s], err: %v", req.URL, err)
		_http.CreateStatusResponse(rw, http.StatusInternalServerError)
		return
	}

	_http.CreateResponse(rw, _http.ContentTypeTextPlain, http.StatusOK, value.String())
}

func (handler *MetricHandler) GetMetricJSON(rw http.ResponseWriter, req *http.Request) {
	var metricReq metric.MetricsDTO

	err := json.NewDecoder(req.Body).Decode(&metricReq)
	if err != nil {
		_http.CreateStatusResponse(rw, http.StatusBadRequest)
		return
	}

	err = handler.checkMetricsCommon(metricReq.MType, metricReq.ID)
	if err != nil {
		handler.logger.Errorf("Incorrect get metric request [%s], JSON: [%v], err: %v", req.URL, metricReq, err)
		CreateResponseOnRequestError(rw, err)
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
			_http.CreateStatusResponse(rw, http.StatusNotFound)
			return
		}

		handler.logger.Errorf("Get metric error on request [%s], JSON: [%v] err: %v", req.URL, metricReq, err)
		_http.CreateStatusResponse(rw, http.StatusInternalServerError)
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

	_http.CreateJSONResponse(rw, http.StatusOK, metricResp)
}

func (handler *MetricHandler) GetMetrics(rw http.ResponseWriter, req *http.Request) {
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
		_http.CreateStatusResponse(rw, http.StatusInternalServerError)
		return
	}
	tmpl, _ := template.New("metrics").Parse(pageTemplate)
	buf := new(bytes.Buffer)
	tmpl.Execute(buf, metrics)
	_http.CreateResponse(rw, _http.ContentTypeHTML, http.StatusOK, buf.String())
}
