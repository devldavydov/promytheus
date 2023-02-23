package storage

import (
	"github.com/devldavydov/promytheus/internal/common/metric"
)

type StorageItem struct {
	MetricName string
	Value      metric.MetricValue
}

type Storage interface {
	SetGaugeMetric(metricName string, value metric.Gauge) (metric.Gauge, error)
	GetGaugeMetric(metricName string) (metric.Gauge, error)
	SetCounterMetric(metricName string, value metric.Counter) (metric.Counter, error)
	GetCounterMetric(metricName string) (metric.Counter, error)
	SetMetrics(metricList []StorageItem) ([]StorageItem, error)
	GetAllMetrics() ([]StorageItem, error)
	Ping() bool
	Close()
}
