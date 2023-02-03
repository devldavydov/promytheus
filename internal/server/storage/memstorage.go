package storage

import (
	"sort"
	"sync"

	"github.com/devldavydov/promytheus/internal/common/metric"
)

type MemStorage struct {
	mu             sync.Mutex
	gaugeStorage   map[string]metric.Gauge
	counterStorage map[string]metric.Counter
}

func NewMemStorage() *MemStorage {
	return &MemStorage{gaugeStorage: make(map[string]metric.Gauge), counterStorage: make(map[string]metric.Counter)}
}

func (storage *MemStorage) SetGaugeMetric(metricName string, value metric.Gauge) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	storage.gaugeStorage[metricName] = value
	return nil
}

func (storage *MemStorage) SetAndGetGaugeMetric(metricName string, value metric.Gauge) (metric.Gauge, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	storage.gaugeStorage[metricName] = value
	return value, nil
}

func (storage *MemStorage) GetGaugeMetric(metricName string) (metric.Gauge, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	val, ok := storage.gaugeStorage[metricName]
	if !ok {
		return 0, ErrMetricNotFound
	}
	return val, nil
}

func (storage *MemStorage) SetCounterMetric(metricName string, value metric.Counter) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	storage.counterStorage[metricName] += value
	return nil
}

func (storage *MemStorage) SetAndGetCounterMetric(metricName string, value metric.Counter) (metric.Counter, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	storage.counterStorage[metricName] += value
	return storage.counterStorage[metricName], nil
}

func (storage *MemStorage) GetCounterMetric(metricName string) (metric.Counter, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	val, ok := storage.counterStorage[metricName]
	if !ok {
		return 0, ErrMetricNotFound
	}
	return val, nil
}

func (storage *MemStorage) GetAllMetrics() ([]StorageItem, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	items := make([]StorageItem, 0, len(storage.counterStorage)+len(storage.gaugeStorage))

	counterItems := sortItems(mapToItems(storage.counterStorage))
	gaugeItems := sortItems(mapToItems(storage.gaugeStorage))

	return append(append(items, counterItems...), gaugeItems...), nil
}

func mapToItems[V metric.MetricValue](m map[string]V) []StorageItem {
	result := make([]StorageItem, 0, len(m))
	for k, v := range m {
		result = append(result, StorageItem{k, v})
	}
	return result
}

func sortItems(items []StorageItem) []StorageItem {
	sort.Slice(items, func(i, j int) bool {
		return items[i].MetricName < items[j].MetricName
	})
	return items
}
