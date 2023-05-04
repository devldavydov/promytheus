package metric

import (
	"fmt"

	"github.com/devldavydov/promytheus/internal/common/hash"
	"github.com/devldavydov/promytheus/internal/common/metric"
)

type requestParams struct {
	metricType   string
	metricName   string
	gaugeValue   metric.Gauge
	counterValue metric.Counter
}

func (handler *MetricHandler) checkMetricsCommon(metricType, metricName string) error {
	if !metric.AllTypes[metricType] {
		return ErrUnknownMetricType
	}

	if len(metricName) == 0 {
		return ErrEmptyMetricName
	}

	return nil
}

func (handler *MetricHandler) parseUpdateRequest(metricType, metricName, metricValue string) (*requestParams, error) {
	err := handler.checkMetricsCommon(metricType, metricName)
	if err != nil {
		return nil, err
	}

	// gauge
	if metric.GaugeTypeName == metricType {
		gaugeVal, err := metric.NewGaugeFromString(metricValue)
		if err != nil {
			return nil, fmt.Errorf("incorrect %s: %w", metric.GaugeTypeName, ErrWrongMetricValue)
		}
		return &requestParams{
			metricType: metric.GaugeTypeName,
			metricName: metricName,
			gaugeValue: gaugeVal,
		}, nil
	}

	// counter
	counterVal, err := metric.NewCounterFromString(metricValue)
	if err != nil {
		return nil, fmt.Errorf("incorrect %s: %w", metric.CounterTypeName, ErrWrongMetricValue)
	}
	return &requestParams{
		metricType:   metric.CounterTypeName,
		metricName:   metricName,
		counterValue: counterVal,
	}, nil
}

func (handler *MetricHandler) parseUpdateRequestJSON(metricReq metric.MetricsDTO) (*requestParams, error) {
	err := handler.checkMetricsCommon(metricReq.MType, metricReq.ID)
	if err != nil {
		return nil, err
	}

	// gauge
	if metric.GaugeTypeName == metricReq.MType {
		gaugeVal, err := metric.NewGaugeFromFloatP(metricReq.Value)
		if err != nil {
			return nil, fmt.Errorf("incorrect %s: %w", metric.GaugeTypeName, ErrWrongMetricValue)
		}

		if err = handler.hmacCheck(metricReq, gaugeVal); err != nil {
			return nil, fmt.Errorf("incorrect %s: %w", metric.GaugeTypeName, err)
		}

		return &requestParams{
			metricType: metric.GaugeTypeName,
			metricName: metricReq.ID,
			gaugeValue: gaugeVal,
		}, nil
	}

	// counter
	counterVal, err := metric.NewCounterFromIntP(metricReq.Delta)
	if err != nil {
		return nil, fmt.Errorf("incorrect %s: %w", metric.CounterTypeName, ErrWrongMetricValue)
	}

	if err = handler.hmacCheck(metricReq, counterVal); err != nil {
		return nil, fmt.Errorf("incorrect %s: %w", metric.CounterTypeName, err)
	}

	return &requestParams{
		metricType:   metric.CounterTypeName,
		metricName:   metricReq.ID,
		counterValue: counterVal,
	}, nil
}

func (handler *MetricHandler) parseUpdateRequestJSONBatch(metricReqList []metric.MetricsDTO) ([]requestParams, error) {
	requestParamsList := make([]requestParams, 0, len(metricReqList))

	for _, metricReq := range metricReqList {
		requestParams, err := handler.parseUpdateRequestJSON(metricReq)
		if err != nil {
			return nil, err
		}
		requestParamsList = append(requestParamsList, *requestParams)
	}

	return requestParamsList, nil
}

func (handler *MetricHandler) hmacCheck(metricReq metric.MetricsDTO, value metric.MetricValue) error {
	if handler.hmacKey == nil {
		return nil
	}

	if !hash.HmacEqual(*metricReq.Hash, value.Hmac(metricReq.ID, *handler.hmacKey)) {
		return ErrMetricHashCheck
	}
	return nil
}
