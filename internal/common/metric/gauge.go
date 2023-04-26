package metric

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/devldavydov/promytheus/internal/common/hash"
)

const GaugeTypeName string = "gauge"

// Gauge - gauge type.
type Gauge float64

var _ MetricValue = (*Gauge)(nil)

// NewGaugeFromString returns new Gauge from string or error.
func NewGaugeFromString(val string) (Gauge, error) {
	flVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}
	return NewGaugeFromFloatP(&flVal)
}

// NewGaugeFromString returns new Gauge from float64 pointer or error.
func NewGaugeFromFloatP(val *float64) (Gauge, error) {
	if val == nil {
		return 0, errors.New("nil pointer")
	}

	return Gauge(*val), nil
}

// FloatP returns Gauge value converted to floa64 pointer.
func (g Gauge) FloatP() *float64 {
	v := float64(g)
	return &v
}

// String returns string representation of Gauge value.
func (g Gauge) String() string {
	return strconv.FormatFloat(float64(g), 'f', 3, 64)
}

// TypeName returns Gauge type name.
func (g Gauge) TypeName() string {
	return GaugeTypeName
}

// Hmac returns hmac value for Gauge.
func (g Gauge) Hmac(id, key string) string {
	return hash.HmacSHA256(fmt.Sprintf("%s:gauge:%f", id, g), key)
}
