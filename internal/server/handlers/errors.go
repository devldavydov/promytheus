package handlers

import "errors"

var ErrorUnknownMetricType = errors.New("unknowm metric type")
var ErrorEmptyMetricName = errors.New("empty metric name")
var ErrorWrongMetricValue = errors.New("wrong metric value")
