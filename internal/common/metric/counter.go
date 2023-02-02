package metric

import (
	"errors"
	"strconv"
)

const CounterTypeName string = "counter"

type Counter int64

func NewCounterFromString(val string) (Counter, error) {
	intVal, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return NewCounterFromIntP(&intVal)
}

func NewCounterFromIntP(val *int64) (Counter, error) {
	if val == nil {
		return 0, errors.New("nil pointer")
	}

	if *val < 0 {
		return 0, errors.New("value below zero")
	}
	return Counter(*val), nil
}

func (c Counter) IntP() *int64 {
	v := int64(c)
	return &v
}

func (c Counter) String() string {
	return strconv.FormatInt(int64(c), 10)
}

func (c Counter) TypeName() string {
	return CounterTypeName
}