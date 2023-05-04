package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/google/uuid"
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
}

func TestCounterSetExistingAndGet(t *testing.T) {
	storage := createMemStorageWithoutPersist()
	storage.SetCounterMetric("foo", metric.Counter(5))
	storage.SetCounterMetric("foo", metric.Counter(5))

	res, err := storage.GetCounterMetric("foo")
	assert.NoError(t, err)
	assert.Equal(t, metric.Counter(10), res)

	res, err = storage.SetCounterMetric("foo", metric.Counter(5))
	assert.NoError(t, err)
	assert.Equal(t, metric.Counter(15), res)
}

func TestCounterGetUnknown(t *testing.T) {
	storage := createMemStorageWithoutPersist()

	_, err := storage.GetCounterMetric("foo")
	assert.ErrorIs(t, err, ErrMetricNotFound)
}

func TestSetMetrics(t *testing.T) {
	storage := createMemStorageWithoutPersist()

	err := storage.SetMetrics([]StorageItem{
		{MetricName: "foo", Value: metric.Gauge(10.1)},
		{MetricName: "cnt1", Value: metric.Counter(1)},
		{MetricName: "cnt1", Value: metric.Counter(1)},
		{MetricName: "cnt1", Value: metric.Counter(1)},
		{MetricName: "cnt2", Value: metric.Counter(2)},
	})
	assert.NoError(t, err)

	result, err := storage.GetAllMetrics()
	assert.NoError(t, err)

	assert.Equal(t, result, []StorageItem{
		{MetricName: "cnt1", Value: metric.Counter(3)},
		{MetricName: "cnt2", Value: metric.Counter(2)},
		{MetricName: "foo", Value: metric.Gauge(10.1)},
	}, result)
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
		{MetricName: "bar", Value: metric.Counter(10)},
		{MetricName: "foo", Value: metric.Counter(5)},
		{MetricName: "buzz", Value: metric.Gauge(1.23456)},
		{MetricName: "fuzz", Value: metric.Gauge(0)},
	}, items)
}

func TestSyncPersistAndRestore(t *testing.T) {
	tmpFile, err := os.CreateTemp("/tmp", "test")
	assert.NoError(t, err)
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	logger := logrus.New()
	storage, err := NewMemStorage(context.TODO(), logger, NewPersistSettings(0, tmpFile.Name(), false))
	assert.NoError(t, err)

	storage.SetCounterMetric("foo", metric.Counter(5))
	storage.SetGaugeMetric("bar", metric.Gauge(4.9))

	storage2, err := NewMemStorage(context.TODO(), logger, NewPersistSettings(0, tmpFile.Name(), true))
	assert.NoError(t, err)

	cVal, err := storage2.GetCounterMetric("foo")
	assert.NoError(t, err)
	assert.Equal(t, metric.Counter(5), cVal)

	gVal, err := storage2.GetGaugeMetric("bar")
	assert.NoError(t, err)
	assert.Equal(t, metric.Gauge(4.9), gVal)
}

func TestSyncIntervalPersistAndRestore(t *testing.T) {
	tmpFile, err := os.CreateTemp("/tmp", "test")
	assert.NoError(t, err)
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	ctx, cancel := context.WithCancel(context.Background())

	logger := logrus.New()
	storage, err := NewMemStorage(ctx, logger, NewPersistSettings(2, tmpFile.Name(), false))
	assert.NoError(t, err)

	storage.SetCounterMetric("foo", metric.Counter(5))
	storage.SetGaugeMetric("bar", metric.Gauge(4.9))

	time.Sleep(5 * time.Second)
	cancel()

	storage2, err := NewMemStorage(context.TODO(), logger, NewPersistSettings(0, tmpFile.Name(), true))
	assert.NoError(t, err)

	cVal, err := storage2.GetCounterMetric("foo")
	assert.NoError(t, err)
	assert.Equal(t, metric.Counter(5), cVal)

	gVal, err := storage2.GetGaugeMetric("bar")
	assert.NoError(t, err)
	assert.Equal(t, metric.Gauge(4.9), gVal)
}

func TestRestoreFileOpenError(t *testing.T) {
	logger := logrus.New()

	_, err := NewMemStorage(context.TODO(), logger, NewPersistSettings(0, "/etc/shadow", true))
	assert.Error(t, err)
}

func TestRestoreFileNotExists(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := logrus.New()

	_, err := NewMemStorage(ctx, logger, NewPersistSettings(0, uuid.NewString(), true))
	assert.NoError(t, err)
}

func TestRestoreFileFormatError(t *testing.T) {
	logger := logrus.New()

	for _, tt := range []struct {
		name     string
		fileData string
	}{
		{name: "wrong json format", fileData: "foobar"},
		{name: "wrong counter format", fileData: `[{"id":"foo","type":"counter","delta":-123}]`},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("/tmp", "test")
			assert.NoError(t, err)
			tmpFile.Write([]byte(tt.fileData))
			tmpFile.Close()
			defer os.Remove(tmpFile.Name())

			_, err = NewMemStorage(context.TODO(), logger, NewPersistSettings(0, tmpFile.Name(), true))
			assert.Error(t, err)
		})
	}

}

func TestPing(t *testing.T) {
	storage := createMemStorageWithoutPersist()
	assert.True(t, storage.Ping())
}

func createMemStorageWithoutPersist() *MemStorage {
	logger := logrus.New()
	storage, _ := NewMemStorage(context.TODO(), logger, NewPersistSettings(0, "", false))
	defer storage.Close()
	return storage
}
