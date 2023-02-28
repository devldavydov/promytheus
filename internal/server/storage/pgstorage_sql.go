package storage

const (
	_sqlCreateTable = `
	CREATE TABLE IF NOT EXISTS metric (
		id    text NOT NULL,
		mtype text NOT NULL,
		delta bigint,
		value double precision,
		
		PRIMARY KEY (id, mtype),
		CHECK(mtype IN ('counter', 'gauge')),
		CHECK(mtype = 'counter' AND delta IS NOT NULL OR mtype = 'gauge' AND value IS NOT NULL)
	);
	`
	_sqlUpsertGauge = `
	INSERT INTO metric (id, mtype, value)
	VALUES ($1, $2, $3)
	ON CONFLICT (id, mtype) DO UPDATE
	SET value = $3
	RETURNING value
	`
	_sqlSelectGauge = `
	SELECT value FROM metric
	WHERE id=$1 AND mtype=$2
	`
	_sqlUpsertCounter = `
	INSERT INTO metric (id, mtype, delta)
	VALUES ($1, $2, $3)
	ON CONFLICT (id, mtype) DO UPDATE
	SET delta = metric.delta + $3
	RETURNING delta
	`
	_sqlSelectCounter = `
	SELECT delta FROM metric
	WHERE id=$1 AND mtype=$2
	`
	_sqlSelectAllMetrics = `
	SELECT id, mtype, delta, value
	FROM metric
	ORDER BY mtype, id
	`
)
