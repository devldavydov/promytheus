package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricsToItemList(t *testing.T) {
	expected := []MetricsItem{
		{"Alloc", "gauge", "0.000"},
		{"BuckHashSys", "gauge", "0.000"},
		{"Frees", "gauge", "0.000"},
		{"GCCPUFraction", "gauge", "0.000"},
		{"GCSys", "gauge", "0.000"},
		{"HeapAlloc", "gauge", "0.000"},
		{"HeapIdle", "gauge", "0.000"},
		{"HeapInuse", "gauge", "0.000"},
		{"HeapObjects", "gauge", "0.000"},
		{"HeapReleased", "gauge", "0.000"},
		{"HeapSys", "gauge", "0.000"},
		{"LastGC", "gauge", "0.000"},
		{"Lookups", "gauge", "0.000"},
		{"MCacheInuse", "gauge", "0.000"},
		{"MCacheSys", "gauge", "0.000"},
		{"MSpanInuse", "gauge", "0.000"},
		{"MSpanSys", "gauge", "0.000"},
		{"Mallocs", "gauge", "0.000"},
		{"NextGC", "gauge", "0.000"},
		{"NumForcedGC", "gauge", "0.000"},
		{"NumGC", "gauge", "0.000"},
		{"OtherSys", "gauge", "0.000"},
		{"PauseTotalNs", "gauge", "0.000"},
		{"StackInuse", "gauge", "0.000"},
		{"StackSys", "gauge", "0.000"},
		{"Sys", "gauge", "0.000"},
		{"TotalAlloc", "gauge", "0.000"},
		{"PollCount", "counter", "0"},
		{"RandomValue", "gauge", "0.000"},
	}
	assert.Equal(t, expected, Metrics{}.ToItemsList())
}
