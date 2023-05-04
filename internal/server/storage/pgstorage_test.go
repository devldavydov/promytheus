package storage

import (
	"context"
	"testing"

	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	_ "github.com/lib/pq"
)

var logger = logrus.New()

const databaseDsn = "postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable"

type PgStorageSuite struct {
	suite.Suite
	stg *PgStorage
}

func (pg *PgStorageSuite) SetupTest() {
	var err error
	pg.stg, err = NewPgStorage(databaseDsn, logger)
	require.NoError(pg.T(), err)
}

func (pg *PgStorageSuite) TearDownTest() {
	pg.stg.Close()
}

func (pg *PgStorageSuite) TestPing() {
	pg.True(pg.stg.Ping())
}

func (pg *PgStorageSuite) TestSetGaugeMetric() {
	metricName := uuid.NewString()

	pg.Run("set first value", func() {
		val, err := pg.stg.SetGaugeMetric(metricName, metric.Gauge(123.123))
		pg.NoError(err)
		pg.Equal(metric.Gauge(123.123), val)
	})

	pg.Run("set second value", func() {
		val, err := pg.stg.SetGaugeMetric(metricName, metric.Gauge(456))
		pg.NoError(err)
		pg.Equal(metric.Gauge(456), val)
	})
}

func (pg *PgStorageSuite) TestGetGaugeMetric() {
	metricName := uuid.NewString()

	pg.Run("get metric - error not found", func() {
		_, err := pg.stg.GetGaugeMetric(metricName)
		pg.ErrorIs(err, ErrMetricNotFound)
	})

	pg.Run("set and get metric", func() {
		_, err := pg.stg.SetGaugeMetric(metricName, metric.Gauge(123.123))
		pg.NoError(err)

		val, err := pg.stg.GetGaugeMetric(metricName)
		pg.NoError(err)
		pg.Equal(metric.Gauge(123.123), val)
	})
}

func (pg *PgStorageSuite) TestSetCounterMetric() {
	metricName := uuid.NewString()

	pg.Run("set first value", func() {
		val, err := pg.stg.SetCounterMetric(metricName, metric.Counter(1))
		pg.NoError(err)
		pg.Equal(metric.Counter(1), val)
	})

	pg.Run("set second value", func() {
		val, err := pg.stg.SetCounterMetric(metricName, metric.Counter(9))
		pg.NoError(err)
		pg.Equal(metric.Counter(10), val)
	})
}

func (pg *PgStorageSuite) TestGetCounterMetric() {
	metricName := uuid.NewString()

	pg.Run("get metric - error not found", func() {
		_, err := pg.stg.GetCounterMetric(metricName)
		pg.ErrorIs(err, ErrMetricNotFound)
	})

	pg.Run("set and get metric", func() {
		_, err := pg.stg.SetCounterMetric(metricName, metric.Counter(123))
		pg.NoError(err)

		val, err := pg.stg.GetCounterMetric(metricName)
		pg.NoError(err)
		pg.Equal(metric.Counter(123), val)
	})
}

func (pg *PgStorageSuite) TestBatchSetAndGet() {
	gaugeMetric, counterMetric := uuid.NewString(), uuid.NewString()

	pg.Run("empty table", func() {
		ctx, cancel := context.WithTimeout(context.Background(), _databaseRequestTimeout)
		defer cancel()

		_, err := pg.stg.db.ExecContext(ctx, "DELETE FROM metric")
		require.NoError(pg.T(), err)
	})

	pg.Run("get all - empty list", func() {
		lst, err := pg.stg.GetAllMetrics()
		pg.NoError(err)
		pg.Equal(0, len(lst))
	})

	pg.Run("set metrics", func() {
		err := pg.stg.SetMetrics([]StorageItem{
			{MetricName: gaugeMetric, Value: metric.Gauge(1.0)},
			{MetricName: gaugeMetric, Value: metric.Gauge(2.0)},
			{MetricName: counterMetric, Value: metric.Counter(1)},
			{MetricName: counterMetric, Value: metric.Counter(2)},
		})
		pg.NoError(err)
	})

	pg.Run("get all", func() {
		lst, err := pg.stg.GetAllMetrics()
		pg.NoError(err)
		pg.Equal(2, len(lst))

		for _, item := range lst {
			if item.MetricName == gaugeMetric {
				pg.Equal(metric.Gauge(2.0), item.Value)
			} else {
				pg.Equal(metric.Counter(3), item.Value)
			}
		}
	})
}

func TestPgStorageSuite(t *testing.T) {
	suite.Run(t, new(PgStorageSuite))
}

func TestPgStorageCreateError(t *testing.T) {
	_, err := NewPgStorage("FooBar", logger)
	assert.Error(t, err)
}
