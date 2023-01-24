package metric

import "fmt"

var AllTypes = map[string]bool{
	GaugeTypeName:   true,
	CounterTypeName: true,
}

type MetricValue interface {
	fmt.Stringer
	TypeName() string
}

type Metrics map[string]MetricValue
