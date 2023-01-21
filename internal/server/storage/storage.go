package storage

import "github.com/devldavydov/promytheus/internal/common/types"

type StorageItem struct {
	MetricName string
	Value      types.MetricValue
}

type Storage interface {
	SetGaugeMetric(metricName string, value types.Gauge) error
	GetGaugeMetric(metricName string) (types.Gauge, error)
	SetCounterMetric(metricName string, value types.Counter) error
	GetCounterMetric(metricName string) (types.Counter, error)
	GetAllMetrics() ([]StorageItem, error)
}
