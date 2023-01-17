package storage

import "github.com/devldavydov/promytheus/internal/common/types"

type MemStorage struct {
	gaugeStorage   map[string]types.Gauge
	counterStorage map[string]types.Counter
}

func NewMemStorage() *MemStorage {
	return &MemStorage{gaugeStorage: make(map[string]types.Gauge), counterStorage: make(map[string]types.Counter)}
}

func (storage *MemStorage) SetGaugeMetric(metricName string, value types.Gauge) error {
	storage.gaugeStorage[metricName] = value
	return nil
}

func (storage *MemStorage) GetGaugeMetric(metricName string) (types.Gauge, error) {
	val, ok := storage.gaugeStorage[metricName]
	if !ok {
		return 0, &MetricNotFoundError{}
	}
	return val, nil
}

func (storage *MemStorage) SetCounterMetric(metricName string, value types.Counter) error {
	storage.counterStorage[metricName] += value
	return nil
}

func (storage *MemStorage) GetCounterMetric(metricName string) (types.Counter, error) {
	val, ok := storage.counterStorage[metricName]
	if !ok {
		return 0, &MetricNotFoundError{}
	}
	return val, nil
}
