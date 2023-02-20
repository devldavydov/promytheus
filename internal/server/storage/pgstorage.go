package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/sirupsen/logrus"
)

const _databaseRequestTimeout = 5 * time.Second

type PgStorage struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewPgStorage(pgConnString string, logger *logrus.Logger) (*PgStorage, error) {
	db, err := sql.Open("postgres", pgConnString)
	if err != nil {
		return nil, err
	}

	pgstorage := &PgStorage{db: db, logger: logger}

	if err = pgstorage.init(); err != nil {
		return nil, err
	}

	return pgstorage, nil
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
	ctx, cancel := context.WithTimeout(context.Background(), _databaseRequestTimeout)
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

func (pgstorage *PgStorage) init() error {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseRequestTimeout)
	defer cancel()

	_, err := pgstorage.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS metric (
			id    text NOT NULL,
			mtype text NOT NULL,
			delta bigint,
			value double precision
		);
	`)
	if err != nil {
		return err
	}

	return nil
}
