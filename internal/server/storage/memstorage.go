package storage

import (
	"sync"

	"github.com/devldavydov/promytheus/internal/common/types"
)

type MemStorage struct {
	mu             sync.Mutex
	gaugeStorage   map[string]types.Gauge
	counterStorage map[string]types.Counter
}

func NewMemStorage() *MemStorage {
	return &MemStorage{gaugeStorage: make(map[string]types.Gauge), counterStorage: make(map[string]types.Counter)}
}

func (storage *MemStorage) SetGaugeMetric(metricName string, value types.Gauge) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	storage.gaugeStorage[metricName] = value
	return nil
}

func (storage *MemStorage) GetGaugeMetric(metricName string) (types.Gauge, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	val, ok := storage.gaugeStorage[metricName]
	if !ok {
		return 0, &MetricNotFoundError{}
	}
	return val, nil
}

func (storage *MemStorage) SetCounterMetric(metricName string, value types.Counter) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	storage.counterStorage[metricName] += value
	return nil
}

func (storage *MemStorage) GetCounterMetric(metricName string) (types.Counter, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	val, ok := storage.counterStorage[metricName]
	if !ok {
		return 0, &MetricNotFoundError{}
	}
	return val, nil
}
