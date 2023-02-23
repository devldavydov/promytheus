package metric

import (
	"encoding/json"
	"net/http"

	_http "github.com/devldavydov/promytheus/internal/common/http"
	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/devldavydov/promytheus/internal/server/storage"
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
		_, err = handler.storage.SetGaugeMetric(params.metricName, params.gaugeValue)
	} else if metric.CounterTypeName == params.metricType {
		_, err = handler.storage.SetCounterMetric(params.metricName, params.counterValue)
	}

	if err != nil {
		handler.logger.Errorf("Update metric error on request [%s], err: %v", req.URL, err)
		_http.CreateStatusResponse(rw, http.StatusInternalServerError)
		return
	}

	_http.CreateStatusResponse(rw, http.StatusOK)
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
		val, err = handler.storage.SetGaugeMetric(params.metricName, params.gaugeValue)
	} else if metric.CounterTypeName == params.metricType {
		val, err = handler.storage.SetCounterMetric(params.metricName, params.counterValue)
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

func (handler *MetricHandler) UpdateMetricJSONBatch(rw http.ResponseWriter, req *http.Request) {
	var metricReqList []metric.MetricsDTO

	err := json.NewDecoder(req.Body).Decode(&metricReqList)
	if err != nil {
		_http.CreateStatusResponse(rw, http.StatusBadRequest)
		return
	}

	// Parse params
	paramsList, err := handler.parseUpdateRequestJSONBatch(metricReqList)
	if err != nil {
		handler.logger.Errorf("Incorrect update metric request [%s], JSON: [%v] , err: %v", req.URL, metricReqList, err)
		CreateResponseOnRequestError(rw, err)
		return
	}

	// Save in storage
	resultMetrics, err := handler.storage.SetMetrics(handler.convertFromParams(paramsList))
	if err != nil {
		handler.logger.Errorf("Update metric error on request [%s], JSON: [%v], err: %v", req.URL, metricReqList, err)
		_http.CreateStatusResponse(rw, http.StatusInternalServerError)
		return
	}

	// Send JSON response with new values
	_http.CreateJSONResponse(rw, http.StatusOK, handler.convertToResponse(resultMetrics))
}

func (handler *MetricHandler) convertFromParams(items []requestParams) []storage.StorageItem {
	storageItemList := make([]storage.StorageItem, 0, len(items))
	for _, params := range items {
		stgItem := storage.StorageItem{MetricName: params.metricName}
		if params.metricType == metric.CounterTypeName {
			stgItem.Value = params.counterValue
		} else if params.metricType == metric.GaugeTypeName {
			stgItem.Value = params.gaugeValue
		}

		storageItemList = append(storageItemList, stgItem)
	}
	return storageItemList
}

func (handler *MetricHandler) convertToResponse(items []storage.StorageItem) []metric.MetricsDTO {
	responseMetrics := make([]metric.MetricsDTO, 0, len(items))

	for _, item := range items {
		responseItem := metric.MetricsDTO{ID: item.MetricName, MType: item.Value.TypeName()}

		if item.Value.TypeName() == metric.CounterTypeName {
			responseItem.Delta = item.Value.(metric.Counter).IntP()
		} else if item.Value.TypeName() == metric.GaugeTypeName {
			responseItem.Value = item.Value.(metric.Gauge).FloatP()
		}

		if handler.hmacKey != nil {
			hash := item.Value.Hmac(item.MetricName, *handler.hmacKey)
			responseItem.Hash = &hash
		}

		responseMetrics = append(responseMetrics, responseItem)
	}

	return responseMetrics
}
