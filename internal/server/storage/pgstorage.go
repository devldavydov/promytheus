package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/sirupsen/logrus"
)

const _databasePingTimeout = 1 * time.Second

type PgStorage struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewPgStorage(pgConnString string, logger *logrus.Logger) (*PgStorage, error) {
	db, err := sql.Open("postgres", pgConnString)
	if err != nil {
		return nil, err
	}
	return &PgStorage{db: db, logger: logger}, nil
}

var _ Storage = (*PgStorage)(nil)

func (pgstorage *PgStorage) SetGaugeMetric(metricName string, value metric.Gauge) error {
	return nil
}

func (pgstorage *PgStorage) SetAndGetGaugeMetric(metricName string, value metric.Gauge) (metric.Gauge, error) {
	return 0, nil
}

func (pgstorage *PgStorage) GetGaugeMetric(metricName string) (metric.Gauge, error) {
	return 0, nil
}

func (pgstorage *PgStorage) SetCounterMetric(metricName string, value metric.Counter) error {
	return nil
}

func (pgstorage *PgStorage) SetAndGetCounterMetric(metricName string, value metric.Counter) (metric.Counter, error) {
	return 0, nil
}

func (pgstorage *PgStorage) GetCounterMetric(metricName string) (metric.Counter, error) {
	return 0, nil
}

func (pgstorage *PgStorage) GetAllMetrics() ([]StorageItem, error) {
	return nil, nil
}

func (pgstorage *PgStorage) Ping() bool {
	ctx, cancel := context.WithTimeout(context.Background(), _databasePingTimeout)
	defer cancel()

	if err := pgstorage.db.PingContext(ctx); err != nil {
		pgstorage.logger.Errorf("Failed to ping database, err: %v", err)
		return false
	}

	return true
}

func (pgstorage *PgStorage) Close() {
	if pgstorage.db == nil {
		return
	}

	err := pgstorage.db.Close()
	if err != nil {
		pgstorage.logger.Errorf("Database conn close err: %v", err)
	}
}
