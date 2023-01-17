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
	assert.Nil(t, err)
	assert.Equal(t, val, res)
}

func TestGaugeGetUnknown(t *testing.T) {
	storage := NewMemStorage()

	_, err := storage.GetGaugeMetric("foo")
	var notFoundErr *MetricNotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestCounterSetNewAndGet(t *testing.T) {
	storage := NewMemStorage()
	val := types.Counter(5)
	storage.SetCounterMetric("foo", val)

	res, err := storage.GetCounterMetric("foo")
	assert.Nil(t, err)
	assert.Equal(t, val, res)
}

func TestCounterSetExistingAndGet(t *testing.T) {
	storage := NewMemStorage()
	storage.SetCounterMetric("foo", types.Counter(5))
	storage.SetCounterMetric("foo", types.Counter(5))

	res, err := storage.GetCounterMetric("foo")
	assert.Nil(t, err)
	assert.Equal(t, types.Counter(10), res)
}

func TestCounterGetUnknown(t *testing.T) {
	storage := NewMemStorage()

	_, err := storage.GetCounterMetric("foo")
	var notFoundErr *MetricNotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}
