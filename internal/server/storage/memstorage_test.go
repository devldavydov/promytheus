package storage

import (
	"context"
	"testing"

	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGaugeSetAndGet(t *testing.T) {
	storage := createMemStorageWithoutPersist()
	val := metric.Gauge(123.456)
	storage.SetGaugeMetric("foo", val)

	res, err := storage.GetGaugeMetric("foo")
	assert.NoError(t, err)
	assert.Equal(t, val, res)

	res, err = storage.SetAndGetGaugeMetric("bar", val)
	assert.NoError(t, err)
	assert.Equal(t, val, res)
}

func TestGaugeGetUnknown(t *testing.T) {
	storage := createMemStorageWithoutPersist()

	_, err := storage.GetGaugeMetric("foo")
	assert.ErrorIs(t, err, ErrMetricNotFound)
}

func TestCounterSetNewAndGet(t *testing.T) {
	storage := createMemStorageWithoutPersist()
	val := metric.Counter(5)
	storage.SetCounterMetric("foo", val)

	res, err := storage.GetCounterMetric("foo")
	assert.NoError(t, err)
	assert.Equal(t, val, res)

	res, err = storage.SetAndGetCounterMetric("bar", val)
	assert.NoError(t, err)
	assert.Equal(t, val, res)
}

func TestCounterSetExistingAndGet(t *testing.T) {
	storage := createMemStorageWithoutPersist()
	storage.SetCounterMetric("foo", metric.Counter(5))
	storage.SetCounterMetric("foo", metric.Counter(5))

	res, err := storage.GetCounterMetric("foo")
	assert.NoError(t, err)
	assert.Equal(t, metric.Counter(10), res)

	res, err = storage.SetAndGetCounterMetric("foo", metric.Counter(5))
	assert.NoError(t, err)
	assert.Equal(t, metric.Counter(15), res)
}

func TestCounterGetUnknown(t *testing.T) {
	storage := createMemStorageWithoutPersist()

	_, err := storage.GetCounterMetric("foo")
	assert.ErrorIs(t, err, ErrMetricNotFound)
}

func TestGetAllMetrics(t *testing.T) {
	storage := createMemStorageWithoutPersist()
	storage.SetCounterMetric("foo", metric.Counter(5))
	storage.SetCounterMetric("bar", metric.Counter(10))
	storage.SetGaugeMetric("fuzz", metric.Gauge(0))
	storage.SetGaugeMetric("buzz", metric.Gauge(1.23456))

	items, err := storage.GetAllMetrics()
	assert.NoError(t, err)
	assert.Equal(t, []StorageItem{
		{"bar", metric.Counter(10)},
		{"foo", metric.Counter(5)},
		{"buzz", metric.Gauge(1.23456)},
		{"fuzz", metric.Gauge(0)},
	}, items)
}

func createMemStorageWithoutPersist() *MemStorage {
	logger := logrus.New()
	storage, _ := NewMemStorage(context.TODO(), logger, NewPersistSettings(0, "", false))
	return storage
}
