package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricsToItemList(t *testing.T) {
	expected := []MetricsItem{
		{"Alloc", "gauge", "0.000000"},
		{"BuckHashSys", "gauge", "0.000000"},
		{"Frees", "gauge", "0.000000"},
		{"GCCPUFraction", "gauge", "0.000000"},
		{"GCSys", "gauge", "0.000000"},
		{"HeapAlloc", "gauge", "0.000000"},
		{"HeapIdle", "gauge", "0.000000"},
		{"HeapInuse", "gauge", "0.000000"},
		{"HeapObjects", "gauge", "0.000000"},
		{"HeapReleased", "gauge", "0.000000"},
		{"HeapSys", "gauge", "0.000000"},
		{"LastGC", "gauge", "0.000000"},
		{"Lookups", "gauge", "0.000000"},
		{"MCacheInuse", "gauge", "0.000000"},
		{"MCacheSys", "gauge", "0.000000"},
		{"MSpanInuse", "gauge", "0.000000"},
		{"MSpanSys", "gauge", "0.000000"},
		{"Mallocs", "gauge", "0.000000"},
		{"NextGC", "gauge", "0.000000"},
		{"NumForcedGC", "gauge", "0.000000"},
		{"NumGC", "gauge", "0.000000"},
		{"OtherSys", "gauge", "0.000000"},
		{"PauseTotalNs", "gauge", "0.000000"},
		{"StackInuse", "gauge", "0.000000"},
		{"StackSys", "gauge", "0.000000"},
		{"Sys", "gauge", "0.000000"},
		{"TotalAlloc", "gauge", "0.000000"},
		{"PollCount", "counter", "0"},
		{"RandomValue", "gauge", "0.000000"},
	}
	assert.Equal(t, expected, Metrics{}.ToItemsList())
}
