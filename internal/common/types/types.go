package types

import "strconv"

type Gauge float64

func (g Gauge) String() string {
	return strconv.FormatFloat(float64(g), 'f', 6, 64)
}

func (g Gauge) TypeName() string {
	return "gauge"
}

func (g Gauge) StringValue() string {
	return g.String()
}

type Counter int64

func (c Counter) String() string {
	return strconv.FormatInt(int64(c), 10)
}

func (c Counter) TypeName() string {
	return "counter"
}

func (c Counter) StringValue() string {
	return c.String()
}
