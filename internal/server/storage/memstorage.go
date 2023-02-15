package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/sirupsen/logrus"
)

type MemStorage struct {
	mu              sync.RWMutex
	persistSettings PersistSettings
	gaugeStorage    map[string]metric.Gauge
	counterStorage  map[string]metric.Counter
	logger          *logrus.Logger
}

var _ Storage = (*MemStorage)(nil)

func NewMemStorage(ctx context.Context, logger *logrus.Logger, persistSettings PersistSettings) (*MemStorage, error) {
	memStorage := &MemStorage{
		persistSettings: persistSettings,
		gaugeStorage:    make(map[string]metric.Gauge),
		counterStorage:  make(map[string]metric.Counter),
		logger:          logger}

	err := memStorage.init(ctx)
	if err != nil {
		return nil, err
	}
	return memStorage, nil
}

func (storage *MemStorage) SetGaugeMetric(metricName string, value metric.Gauge) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	storage.gaugeStorage[metricName] = value
	storage.trySyncPersist()

	return nil
}

func (storage *MemStorage) SetAndGetGaugeMetric(metricName string, value metric.Gauge) (metric.Gauge, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	storage.gaugeStorage[metricName] = value
	storage.trySyncPersist()

	return value, nil
}

func (storage *MemStorage) GetGaugeMetric(metricName string) (metric.Gauge, error) {
	storage.mu.RLock()
	defer storage.mu.RUnlock()

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
	storage.trySyncPersist()

	return nil
}

func (storage *MemStorage) SetAndGetCounterMetric(metricName string, value metric.Counter) (metric.Counter, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	storage.counterStorage[metricName] += value
	storage.trySyncPersist()

	return storage.counterStorage[metricName], nil
}

func (storage *MemStorage) GetCounterMetric(metricName string) (metric.Counter, error) {
	storage.mu.RLock()
	defer storage.mu.RUnlock()

	val, ok := storage.counterStorage[metricName]
	if !ok {
		return 0, ErrMetricNotFound
	}
	return val, nil
}

func (storage *MemStorage) GetAllMetrics() ([]StorageItem, error) {
	storage.mu.RLock()
	defer storage.mu.RUnlock()

	items := make([]StorageItem, 0, len(storage.counterStorage)+len(storage.gaugeStorage))

	counterItems := sortItems(mapToItems(storage.counterStorage))
	gaugeItems := sortItems(mapToItems(storage.gaugeStorage))

	return append(append(items, counterItems...), gaugeItems...), nil
}

func (storage *MemStorage) init(ctx context.Context) error {
	if storage.persistSettings.ShouldRestore() {
		err := storage.restore()
		if err != nil {
			return err
		}
	}

	if !storage.persistSettings.ShouldPersist() {
		storage.logger.Warning("Storage persist not enabled")
		return nil
	}

	if storage.persistSettings.ShouldIntervalPersist() {
		go storage.persistIntervalThread(ctx)
	}

	return nil
}

func (storage *MemStorage) restore() error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	file, err := os.OpenFile(storage.persistSettings.StoreFile, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			storage.logger.Warnf("Restore file [%s] not exists, skipping...", storage.persistSettings.StoreFile)
			return nil
		} else {
			return fmt.Errorf("restore open file [%s] err: %w", storage.persistSettings.StoreFile, err)
		}
	}
	defer file.Close()

	var totalMetrics []metric.MetricsDTO
	err = json.NewDecoder(file).Decode(&totalMetrics)
	if err != nil {
		return fmt.Errorf("failed to restore storage: %w", err)
	}

	restoreErr := func(v metric.MetricsDTO, err error) error {
		return fmt.Errorf("failed to restore value [%s] of type [%s]: %w", v.ID, v.MType, err)
	}

	for _, v := range totalMetrics {
		if metric.CounterTypeName == v.MType {
			val, err := metric.NewCounterFromIntP(v.Delta)
			if err != nil {
				return restoreErr(v, err)
			}
			storage.counterStorage[v.ID] = val
		} else if metric.GaugeTypeName == v.MType {
			val, err := metric.NewGaugeFromFloatP(v.Value)
			if err != nil {
				return restoreErr(v, err)
			}
			storage.gaugeStorage[v.ID] = val
		}
	}
	storage.logger.Infof("Storage restored from file [%s]", storage.persistSettings.StoreFile)

	return nil
}

func (storage *MemStorage) trySyncPersist() {
	if storage.persistSettings.ShouldSyncPersist() {
		storage.persist()
	}
}

func (storage *MemStorage) persist() {
	totalMetrics := make([]metric.MetricsDTO, 0, len(storage.counterStorage)+len(storage.gaugeStorage))
	for k, v := range storage.counterStorage {
		totalMetrics = append(totalMetrics, metric.MetricsDTO{ID: k, MType: metric.CounterTypeName, Delta: v.IntP()})
	}
	for k, v := range storage.gaugeStorage {
		totalMetrics = append(totalMetrics, metric.MetricsDTO{ID: k, MType: metric.GaugeTypeName, Value: v.FloatP()})
	}

	file, err := os.OpenFile(storage.persistSettings.StoreFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		storage.logger.Errorf("Failed to open persist storage file [%s]: %v", storage.persistSettings.StoreFile, err)
		return
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(totalMetrics)
	if err != nil {
		storage.logger.Errorf("Failed to persist storage: %v", err)
	}
}

func (storage *MemStorage) persistIntervalThread(ctx context.Context) {
	ticker := time.NewTicker(storage.persistSettings.StoreInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			storage.mu.RLock()
			storage.persist()
			storage.mu.RUnlock()
		case <-ctx.Done():
			storage.logger.Info("Storage persist interval thread context canceled")
			return
		}
	}
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
