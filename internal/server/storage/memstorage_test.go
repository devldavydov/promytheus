package storage

import (
	"errors"
	"testing"

	"github.com/devldavydov/promytheus/internal/common/types"
	"github.com/stretchr/testify/assert"
)

func TestGaugeSetAndGet(t *testing.T) {
	storage := NewMemStorage()
	val := types.Gauge(123.456)
	storage.SetGaugeMetric("foo", val)

	res, err := storage.GetGaugeMetric("foo")
	assert.NoError(t, err)
	assert.Equal(t, val, res)
}

func TestGaugeGetUnknown(t *testing.T) {
	storage := NewMemStorage()

	_, err := storage.GetGaugeMetric("foo")
	assert.True(t, errors.As(err, &MetricNotFoundErrorP))
}

func TestCounterSetNewAndGet(t *testing.T) {
	storage := NewMemStorage()
	val := types.Counter(5)
	storage.SetCounterMetric("foo", val)

	res, err := storage.GetCounterMetric("foo")
	assert.NoError(t, err)
	assert.Equal(t, val, res)
}

func TestCounterSetExistingAndGet(t *testing.T) {
	storage := NewMemStorage()
	storage.SetCounterMetric("foo", types.Counter(5))
	storage.SetCounterMetric("foo", types.Counter(5))

	res, err := storage.GetCounterMetric("foo")
	assert.NoError(t, err)
	assert.Equal(t, types.Counter(10), res)
}

func TestCounterGetUnknown(t *testing.T) {
	storage := NewMemStorage()

	_, err := storage.GetCounterMetric("foo")
	var notFoundErr *MetricNotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestGetAllMetrics(t *testing.T) {
	storage := NewMemStorage()
	storage.SetCounterMetric("foo", types.Counter(5))
	storage.SetCounterMetric("bar", types.Counter(10))
	storage.SetGaugeMetric("fuzz", types.Gauge(0))
	storage.SetGaugeMetric("buzz", types.Gauge(1.23456))

	items, err := storage.GetAllMetrics()
	assert.NoError(t, err)
	assert.Equal(t, []StorageItem{
		{"bar", types.Counter(10)},
		{"foo", types.Counter(5)},
		{"buzz", types.Gauge(1.23456)},
		{"fuzz", types.Gauge(0)},
	}, items)
}
