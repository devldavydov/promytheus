package handlers

type IncorrectURLWrongPartsCountError struct {
	err string
}

func (e *IncorrectURLWrongPartsCountError) Error() string {
	return e.err
}

var IncorrectURLWrongPartsCountErrorP *IncorrectURLWrongPartsCountError

type IncorrectURLUnknownMetricType struct {
	err string
}

func (e *IncorrectURLUnknownMetricType) Error() string {
	return e.err
}

var IncorrectURLUnknownMetricTypeP *IncorrectURLUnknownMetricType

type IncorrectURLEmptyMetricName struct {
	err string
}

func (e *IncorrectURLEmptyMetricName) Error() string {
	return e.err
}

var IncorrectURLEmptyMetricNameP *IncorrectURLEmptyMetricName

type IncorrectURLWrongMetricValue struct {
	err string
}

func (e *IncorrectURLWrongMetricValue) Error() string {
	return e.err
}

var IncorrectURLWrongMetricValueP *IncorrectURLWrongMetricValue
