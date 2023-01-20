package types

import "fmt"

var AllTypes = map[string]bool{
	GaugeTypeName:   true,
	CounterTypeName: true,
}

type MetricValue interface {
	fmt.Stringer
	TypeName() string
}
