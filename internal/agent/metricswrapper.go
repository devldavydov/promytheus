package agent

import (
	"sync"

	"github.com/devldavydov/promytheus/internal/agent/metrics"
)

type MetricsWrapper struct {
	metrics metrics.Metrics
	mu      sync.Mutex
}

func (metricsWrapper *MetricsWrapper) Set(metrics metrics.Metrics) {
	metricsWrapper.mu.Lock()
	defer metricsWrapper.mu.Unlock()
	metricsWrapper.metrics = metrics
}

func (metricsWrapper *MetricsWrapper) Get() metrics.Metrics {
	metricsWrapper.mu.Lock()
	defer metricsWrapper.mu.Unlock()
	return metricsWrapper.metrics
}
