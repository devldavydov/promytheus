package handlers

type IncorrectUrlWrongPartsCountError struct {
	err string
}

func (e *IncorrectUrlWrongPartsCountError) Error() string {
	return e.err
}

var IncorrectUrlWrongPartsCountErrorP *IncorrectUrlWrongPartsCountError

type IncorrectUrlUnknownMetricType struct {
	err string
}

func (e *IncorrectUrlUnknownMetricType) Error() string {
	return e.err
}

var IncorrectUrlUnknownMetricTypeP *IncorrectUrlUnknownMetricType

type IncorrectUrlEmptyMetricName struct {
	err string
}

func (e *IncorrectUrlEmptyMetricName) Error() string {
	return e.err
}

var IncorrectUrlEmptyMetricNameP *IncorrectUrlEmptyMetricName

type IncorrectUrlWrongMetricValue struct {
	err string
}

func (e *IncorrectUrlWrongMetricValue) Error() string {
	return e.err
}

var IncorrectUrlWrongMetricValueP *IncorrectUrlWrongMetricValue
