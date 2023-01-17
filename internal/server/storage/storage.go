package storage

import "github.com/devldavydov/promytheus/internal/common/types"

type Storage interface {
	SetGaugeMetric(metricName string, value types.Gauge) error
	GetGaugeMetric(metricName string) (types.Gauge, error)
	SetCounterMetric(metricName string, value types.Counter) error
	GetCounterMetric(metricName string) (types.Counter, error)
}

type MetricNotFoundError struct {
	err string
}

func (e MetricNotFoundError) Error() string {
	return e.err
}
