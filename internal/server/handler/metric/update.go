package metric

import (
	"encoding/json"
	"net/http"

	_http "github.com/devldavydov/promytheus/internal/common/http"
	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

// UpdateMetric set new value for metric.
//
//	@Summary	Update metric
//	@Produce	plain/text
//	@Param		metricType	path	string	true	"Metric Type"
//	@Param		metricName	path	string	true	"Metric Name"
//	@Param		metricValue	path	string	true	"Metric Value"
//	@Success	200			"Updated successfully"
//	@Failure	400			"Bad request"
//	@Failure	500			"Internal error"
//	@Failure	501			"Metric type not found"
//	@Router		/update/{metricType}/{metricName}/{metricValue} [post]
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

// UpdateMetricJSON set new value for metric in JSON.
//
//	@Summary	Update metric in JSON
//	@Accept		json
//	@Produce	json
//	@Param		message	body		metric.MetricsDTO	true	"Metric update request"
//	@Success	200		{object}	metric.MetricsDTO	"Returns updated metric"
//	@Failure	400		"Bad request"
//	@Failure	500		"Internal error"
//	@Failure	501		"Metric type not found"
//	@Router		/update [post]
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

// UpdateMetricJSONBatch set new values for batch of metrics in JSON.
//
//	@Summary	Update metrics batch in JSON
//	@Accept		json
//	@Produce	json
//	@Param		message	body	[]metric.MetricsDTO	true	"Metrics update batch request"
//	@Success	200		{array}	array				"Returns empty array"
//	@Failure	400		"Bad request"
//	@Failure	500		"Internal error"
//	@Failure	501		"Metric type not found"
//	@Router		/updates [post]
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
	if err = handler.storage.SetMetrics(handler.convertFromParams(paramsList)); err != nil {
		handler.logger.Errorf("Update metric error on request [%s], JSON: [%v], err: %v", req.URL, metricReqList, err)
		_http.CreateStatusResponse(rw, http.StatusInternalServerError)
		return
	}

	// Send JSON response with new values
	_http.CreateResponse(rw, _http.ContentTypeApplicationJSON, http.StatusOK, "[]")
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
