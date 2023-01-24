package handler

import "errors"

var ErrUnknownMetricType = errors.New("unknowm metric type")
var ErrEmptyMetricName = errors.New("empty metric name")
var ErrWrongMetricValue = errors.New("wrong metric value")
