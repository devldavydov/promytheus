package storage

import "github.com/devldavydov/promytheus/internal/common/types"

type StorageItem struct {
	MetricName string
	Value      types.BaseType
}

type StorageItemByMetricTypeName []StorageItem

func (s StorageItemByMetricTypeName) Len() int      { return len(s) }
func (s StorageItemByMetricTypeName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s StorageItemByMetricTypeName) Less(i, j int) bool {
	return s[i].Value.TypeName()+s[i].MetricName < s[j].Value.TypeName()+s[j].MetricName
}

type Storage interface {
	SetGaugeMetric(metricName string, value types.Gauge) error
	GetGaugeMetric(metricName string) (types.Gauge, error)
	SetCounterMetric(metricName string, value types.Counter) error
	GetCounterMetric(metricName string) (types.Counter, error)
	GetAllMetrics() ([]StorageItem, error)
}
