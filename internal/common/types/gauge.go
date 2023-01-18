package types

import (
	"strconv"
)

const GaugeTypeName string = "gauge"

type Gauge float64

func NewGaugeFromString(val string) (Gauge, error) {
	flVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}
	return Gauge(flVal), nil
}

func (g Gauge) String() string {
	return strconv.FormatFloat(float64(g), 'f', 3, 64)
}

func (g Gauge) TypeName() string {
	return GaugeTypeName
}
