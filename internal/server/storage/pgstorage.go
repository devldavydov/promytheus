package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/sirupsen/logrus"
)

const _databaseRequestTimeout = 10 * time.Second

type PgStorage struct {
	db     *sql.DB
	logger *logrus.Logger
}

type tableRow struct {
	id    string
	mtype string
	delta sql.NullInt64
	value sql.NullFloat64
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

func (pgstorage *PgStorage) SetGaugeMetric(metricName string, value metric.Gauge) (metric.Gauge, error) {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseRequestTimeout)
	defer cancel()

	var val metric.Gauge
	err := pgstorage.db.QueryRowContext(ctx, _sqlUpsertGauge, metricName, metric.GaugeTypeName, value).Scan(&val)

	if err != nil {
		return 0, err
	}

	return val, nil
}

func (pgstorage *PgStorage) GetGaugeMetric(metricName string) (metric.Gauge, error) {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseRequestTimeout)
	defer cancel()

	var val metric.Gauge
	err := pgstorage.db.QueryRowContext(ctx, _sqlSelectGauge, metricName, metric.GaugeTypeName).Scan(&val)

	switch {
	case err == sql.ErrNoRows:
		return 0, ErrMetricNotFound
	case err != nil:
		return 0, err
	}

	return val, nil
}

func (pgstorage *PgStorage) SetCounterMetric(metricName string, value metric.Counter) (metric.Counter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseRequestTimeout)
	defer cancel()

	var val metric.Counter
	err := pgstorage.db.QueryRowContext(ctx, _sqlUpsertCounter, metricName, metric.CounterTypeName, value).Scan(&val)

	if err != nil {
		return 0, err
	}

	return val, nil
}

func (pgstorage *PgStorage) GetCounterMetric(metricName string) (metric.Counter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseRequestTimeout)
	defer cancel()

	var val metric.Counter
	err := pgstorage.db.QueryRowContext(ctx, _sqlSelectCounter, metricName, metric.CounterTypeName).Scan(&val)

	switch {
	case err == sql.ErrNoRows:
		return 0, ErrMetricNotFound
	case err != nil:
		return 0, err
	}

	return val, nil
}

func (pgstorage *PgStorage) SetMetrics(metricList []StorageItem) ([]StorageItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseRequestTimeout)
	defer cancel()

	// Prepare statements
	tx, err := pgstorage.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmtCounter, err := tx.PrepareContext(ctx, _sqlUpsertCounter)
	if err != nil {
		return nil, err
	}
	defer stmtCounter.Close()

	stmtGauge, err := tx.PrepareContext(ctx, _sqlUpsertGauge)
	if err != nil {
		return nil, err
	}
	defer stmtGauge.Close()

	// Define map for result values
	resultValues := make(map[string]metric.MetricValue)

	// DB operations
	for _, metricItem := range metricList {
		var val metric.MetricValue
		if metricItem.Value.TypeName() == metric.CounterTypeName {
			var valC metric.Counter
			err = stmtCounter.QueryRowContext(ctx,
				metricItem.MetricName,
				metric.CounterTypeName,
				metricItem.Value.(metric.Counter)).Scan(&valC)
			if err != nil {
				return nil, err
			}
			val = valC
		} else if metricItem.Value.TypeName() == metric.GaugeTypeName {
			var valG metric.Gauge
			err = stmtGauge.QueryRowContext(ctx,
				metricItem.MetricName,
				metric.GaugeTypeName,
				metricItem.Value.(metric.Gauge)).Scan(&valG)
			if err != nil {
				return nil, err
			}
			val = valG
		}

		resultValues[metricItem.MetricName] = val
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Prepare result values
	uniqueMetricNames := make(map[string]bool)
	resultStorageItems := make([]StorageItem, 0, len(metricList))

	for _, metricItem := range metricList {
		if uniqueMetricNames[metricItem.MetricName] {
			continue
		}

		uniqueMetricNames[metricItem.MetricName] = true
		resultStorageItems = append(resultStorageItems, StorageItem{MetricName: metricItem.MetricName, Value: resultValues[metricItem.MetricName]})
	}

	return resultStorageItems, nil
}

func (pgstorage *PgStorage) GetAllMetrics() ([]StorageItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseRequestTimeout)
	defer cancel()

	var items []StorageItem

	rows, err := pgstorage.db.QueryContext(ctx, `
		SELECT id, mtype, delta, value
		FROM metric
		ORDER BY mtype, id
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var r tableRow
		err = rows.Scan(&r.id, &r.mtype, &r.delta, &r.value)
		if err != nil {
			return nil, err
		}

		item := StorageItem{MetricName: r.id}

		if r.mtype == metric.CounterTypeName {
			item.Value = metric.Counter(r.delta.Int64)
		} else if r.mtype == metric.GaugeTypeName {
			item.Value = metric.Gauge(r.value.Float64)
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
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

	_, err := pgstorage.db.ExecContext(ctx, _sqlCreateTable)
	if err != nil {
		return err
	}

	return nil
}
