// Package metrics - base package for metric values and their functions.
package metric

import "fmt"

// AllTypes - valid metric types.
var AllTypes = map[string]bool{
	GaugeTypeName:   true,
	CounterTypeName: true,
}

// MetricValue - value common interface.
type MetricValue interface {
	fmt.Stringer

	TypeName() string
	Hmac(id, key string) string
}

// Metrics - represensts map of MetricValue.
type Metrics map[string]MetricValue

// MetricsDTO - metric structure for JSON serde.
type MetricsDTO struct {
	ID    string   `json:"id"`              // metric name
	MType string   `json:"type"`            // metric type - gauge|counter
	Delta *int64   `json:"delta,omitempty"` // metric value if counter
	Value *float64 `json:"value,omitempty"` // metric value if gauge
	Hash  *string  `json:"hash,omitempty"`  // значение хеш-функции
}
