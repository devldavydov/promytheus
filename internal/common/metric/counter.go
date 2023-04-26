package metric

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/devldavydov/promytheus/internal/common/hash"
)

const CounterTypeName string = "counter"

// Counter - counter type.
type Counter int64

var _ MetricValue = (*Counter)(nil)

// NewCounterFromString returns new Counter from string or error.
func NewCounterFromString(val string) (Counter, error) {
	intVal, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return NewCounterFromIntP(&intVal)
}

// NewCounterFromIntP returns new Counter from int64 pointer or error.
func NewCounterFromIntP(val *int64) (Counter, error) {
	if val == nil {
		return 0, errors.New("nil pointer")
	}

	if *val < 0 {
		return 0, errors.New("value below zero")
	}
	return Counter(*val), nil
}

// IntP returns Counter value converted to int64 pointer.
func (c Counter) IntP() *int64 {
	v := int64(c)
	return &v
}

// String returns string representation of Counter value.
func (c Counter) String() string {
	return strconv.FormatInt(int64(c), 10)
}

// TypeName returns Counter type name.
func (c Counter) TypeName() string {
	return CounterTypeName
}

// Hmac returns hmac value for Counter.
func (c Counter) Hmac(id, key string) string {
	return hash.HmacSHA256(fmt.Sprintf("%s:counter:%d", id, c), key)
}
