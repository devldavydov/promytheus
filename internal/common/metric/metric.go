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

type MetricsDTO struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
