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
	return nil
}

func (storage *MemStorage) GetGaugeMetric(metricName string) (types.Gauge, error) {
	return 0, nil
}

func (storage *MemStorage) SetCounterMetric(metricName string, value types.Counter) error {
	return nil
}

func (storage *MemStorage) GetCounterMetric(metricName string) (types.Counter, error) {
	return 0, nil
}
