package types

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
	if intVal < 0 {
		return 0, errors.New("value below zero")
	}
	return Counter(intVal), nil
}

func (c Counter) String() string {
	return strconv.FormatInt(int64(c), 10)
}

func (c Counter) TypeName() string {
	return CounterTypeName
}
