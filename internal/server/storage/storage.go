// Package storage represensts metrics storage functionality.
package storage

import (
	"github.com/devldavydov/promytheus/internal/common/metric"
)

type StorageItem struct {
	MetricName string
	Value      metric.MetricValue
}

// Storage is a interface for metrics store functionality.
type Storage interface {
	// SetGaugeMetric set gauge metric and return new value from storage or error.
	SetGaugeMetric(metricName string, value metric.Gauge) (metric.Gauge, error)
	// GetGaugeMetric get gauge metric from storage or error.
	GetGaugeMetric(metricName string) (metric.Gauge, error)
	// SetCounterMetric set counter metric and return new value from storage or error.
	SetCounterMetric(metricName string, value metric.Counter) (metric.Counter, error)
	// GetCounterMetric get counter metric from storage or error.
	GetCounterMetric(metricName string) (metric.Counter, error)
	// SetMetrics set list of metrics.
	SetMetrics(metricList []StorageItem) error
	// GetAllMetrics returns list of metrics from storage.
	GetAllMetrics() ([]StorageItem, error)
	// Ping tests storage availability.
	Ping() bool
	// Close storage connection.
	Close()
}
