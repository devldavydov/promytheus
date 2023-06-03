package publisher

import "github.com/devldavydov/promytheus/internal/common/metric"

func iterateMetrics(metricsList []metric.Metrics, fn func(name string, value metric.MetricValue)) {
	for _, metrics := range metricsList {
		for name, value := range metrics {
			fn(name, value)
		}
	}
}

func totalMetrics(metricsList []metric.Metrics) int {
	cnt := 0
	for _, m := range metricsList {
		cnt += len(m)
	}
	return cnt
}

func prepareMetric(metricName string, metricValue metric.MetricValue, hmacKey *string) metric.MetricsDTO {
	metricReq := metric.MetricsDTO{
		ID:    metricName,
		MType: metricValue.TypeName(),
	}

	if metric.GaugeTypeName == metricValue.TypeName() {
		metricReq.Value = metricValue.(metric.Gauge).FloatP()
	} else if metric.CounterTypeName == metricValue.TypeName() {
		metricReq.Delta = metricValue.(metric.Counter).IntP()
	}

	if hmacKey != nil {
		hash := metricValue.Hmac(metricName, *hmacKey)
		metricReq.Hash = &hash
	}

	return metricReq
}
