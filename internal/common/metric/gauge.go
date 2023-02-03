package metric

import (
	"errors"
	"strconv"
)

const GaugeTypeName string = "gauge"

type Gauge float64

func NewGaugeFromString(val string) (Gauge, error) {
	flVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}
	return NewGaugeFromFloatP(&flVal)
}

func NewGaugeFromFloatP(val *float64) (Gauge, error) {
	if val == nil {
		return 0, errors.New("nil pointer")
	}

	return Gauge(*val), nil
}

func (g Gauge) FloatP() *float64 {
	v := float64(g)
	return &v
}

func (g Gauge) String() string {
	return strconv.FormatFloat(float64(g), 'f', 3, 64)
}

func (g Gauge) TypeName() string {
	return GaugeTypeName
}
