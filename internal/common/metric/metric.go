// Package metric - base package for metric values and their functions.
package metric

import (
	"errors"
	"fmt"
)

// Metric errors
var (
	ErrUnknownMetricType = errors.New("unknown metric type")
	ErrEmptyMetricName   = errors.New("empty metric name")
	ErrWrongMetricValue  = errors.New("wrong metric value")
	ErrMetricHashCheck   = errors.New("metric hash check fail")
)

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
	Delta *int64   `json:"delta,omitempty"` // metric value if counter
	Value *float64 `json:"value,omitempty"` // metric value if gauge
	Hash  *string  `json:"hash,omitempty"`  // hash value
	ID    string   `json:"id"`              // metric name
	MType string   `json:"type"`            // metric type - gauge|counter
}
