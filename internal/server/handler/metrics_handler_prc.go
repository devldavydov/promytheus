package handler

import (
	"fmt"

	"github.com/devldavydov/promytheus/internal/common/metric"
)

type requestParams struct {
	metricType   string
	metricName   string
	gaugeValue   metric.Gauge
	counterValue metric.Counter
}

func (handler *MetricsHandler) checkMetricsCommon(metricType, metricName string) error {
	if !metric.AllTypes[metricType] {
		return ErrUnknownMetricType
	}

	if len(metricName) == 0 {
		return ErrEmptyMetricName
	}

	return nil
}

func (handler *MetricsHandler) parseUpdateRequest(metricType, metricName, metricValue string) (requestParams, error) {
	err := handler.checkMetricsCommon(metricType, metricName)
	if err != nil {
		return requestParams{}, err
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

	return requestParams{}, ErrUnknownMetricType
}

func (handler *MetricsHandler) parseUpdateRequestJSON(metricReq metric.MetricsDTO) (requestParams, error) {
	err := handler.checkMetricsCommon(metricReq.MType, metricReq.ID)
	if err != nil {
		return requestParams{}, err
	}

	if metric.GaugeTypeName == metricReq.MType {
		gaugeVal, err := metric.NewGaugeFromFloatP(metricReq.Value)
		if err != nil {
			return requestParams{}, fmt.Errorf("incorrect %s: %w", metric.GaugeTypeName, ErrWrongMetricValue)
		}
		return requestParams{
			metricType: metric.GaugeTypeName,
			metricName: metricReq.ID,
			gaugeValue: gaugeVal,
		}, nil
	} else if metric.CounterTypeName == metricReq.MType {
		counterVal, err := metric.NewCounterFromIntP(metricReq.Delta)
		if err != nil {
			return requestParams{}, fmt.Errorf("incorrect %s: %w", metric.CounterTypeName, ErrWrongMetricValue)
		}
		return requestParams{
			metricType:   metric.CounterTypeName,
			metricName:   metricReq.ID,
			counterValue: counterVal,
		}, nil
	}

	return requestParams{}, ErrUnknownMetricType
}
