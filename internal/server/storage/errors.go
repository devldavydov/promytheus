package storage

type MetricNotFoundError struct {
	err string
}

func (e *MetricNotFoundError) Error() string {
	return e.err
}
