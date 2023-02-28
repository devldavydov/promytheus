package metric

import "fmt"

var AllTypes = map[string]bool{
	GaugeTypeName:   true,
	CounterTypeName: true,
}

type MetricValue interface {
	fmt.Stringer

	TypeName() string
	Hmac(id, key string) string
}

type Metrics map[string]MetricValue

type MetricsDTO struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  *string  `json:"hash,omitempty"`  // значение хеш-функции
}
