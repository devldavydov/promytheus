package metric

import (
	"encoding/json"
	"net/http"

	_http "github.com/devldavydov/promytheus/internal/common/http"
	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/go-chi/chi/v5"
)

func (handler *MetricHandler) UpdateMetric(rw http.ResponseWriter, req *http.Request) {
	params, err := handler.parseUpdateRequest(chi.URLParam(req, "metricType"), chi.URLParam(req, "metricName"), chi.URLParam(req, "metricValue"))
	if err != nil {
		handler.logger.Errorf("Incorrect update metric request [%s], err: %v", req.URL, err)
		CreateResponseOnRequestError(rw, err)
		return
	}

	if metric.GaugeTypeName == params.metricType {
		err = handler.storage.SetGaugeMetric(params.metricName, params.gaugeValue)
	} else if metric.CounterTypeName == params.metricType {
		err = handler.storage.SetCounterMetric(params.metricName, params.counterValue)
	}

	if err != nil {
		handler.logger.Errorf("Update metric error on request [%s], err: %v", req.URL, err)
		_http.CreateStatusResponse(rw, http.StatusInternalServerError)
		return
	}

	_http.CreateResponse(rw, _http.ContentTypeTextPlain, http.StatusOK, http.StatusText(http.StatusOK))
}

func (handler *MetricHandler) UpdateMetricJSON(rw http.ResponseWriter, req *http.Request) {
	var metricReq metric.MetricsDTO

	err := json.NewDecoder(req.Body).Decode(&metricReq)
	if err != nil {
		_http.CreateStatusResponse(rw, http.StatusBadRequest)
		return
	}

	params, err := handler.parseUpdateRequestJSON(metricReq)
	if err != nil {
		handler.logger.Errorf("Incorrect update metric request [%s], JSON: [%v] , err: %v", req.URL, metricReq, err)
		CreateResponseOnRequestError(rw, err)
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
		_http.CreateStatusResponse(rw, http.StatusInternalServerError)
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

	_http.CreateJSONResponse(rw, http.StatusOK, metricResp)
}
