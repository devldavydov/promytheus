package metric

import (
	"fmt"
	"testing"

	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/stretchr/testify/assert"
)

func BenchmarkParseUpdateRequest(b *testing.B) {
	handler := &MetricHandler{}

	b.Run("parse plain gauge", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := handler.parseUpdateRequest(metric.GaugeTypeName, "test", "123.123")
			assert.NoError(b, err)
		}
	})

	b.Run("parse plain counter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := handler.parseUpdateRequest(metric.CounterTypeName, "test", "123")
			assert.NoError(b, err)
		}
	})
}

func BenchmarkParseUpdateJSONRequest(b *testing.B) {
	handler := &MetricHandler{}

	b.Run("parse json gauge", func(b *testing.B) {
		var v float64 = 123.123
		req := metric.MetricsDTO{
			ID:    "test",
			MType: metric.GaugeTypeName,
			Value: &v,
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := handler.parseUpdateRequestJSON(req)
			assert.NoError(b, err)
		}
	})

	b.Run("parse json counter", func(b *testing.B) {
		var v int64 = 123
		req := metric.MetricsDTO{
			ID:    "test",
			MType: metric.CounterTypeName,
			Delta: &v,
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := handler.parseUpdateRequestJSON(req)
			assert.NoError(b, err)
		}
	})
}

func BenchmarkParseUpdateRequestJSONBatch(b *testing.B) {
	handler := &MetricHandler{}
	benchReqCount := 1000000

	reqList := make([]metric.MetricsDTO, 0, benchReqCount)
	for i := 0; i < benchReqCount/2; i++ {
		val := float64(i)
		reqList = append(reqList, metric.MetricsDTO{
			ID:    fmt.Sprintf("metric%d", i),
			MType: metric.GaugeTypeName,
			Value: &val,
		})
	}
	for i := benchReqCount / 2; i < benchReqCount; i++ {
		val := int64(i)
		reqList = append(reqList, metric.MetricsDTO{
			ID:    fmt.Sprintf("metric%d", i),
			MType: metric.CounterTypeName,
			Delta: &val,
		})
	}

	b.Run("parse json batch", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := handler.parseUpdateRequestJSONBatch(reqList)
			assert.NoError(b, err)
		}
	})
}
