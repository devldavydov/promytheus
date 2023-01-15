package metrics

import "sync"

type gauge float64
type counter int64

type Metrics struct {
	Alloc         gauge
	BuckHashSys   gauge
	Frees         gauge
	GCCPUFraction gauge
	GCSys         gauge
	HeapAlloc     gauge
	HeapIdle      gauge
	HeapInuse     gauge
	HeapObjects   gauge
	HeapReleased  gauge
	HeapSys       gauge
	LastGC        gauge
	Lookups       gauge
	MCacheInuse   gauge
	MCacheSys     gauge
	MSpanInuse    gauge
	MSpanSys      gauge
	Mallocs       gauge
	NextGC        gauge
	NumForcedGC   gauge
	NumGC         gauge
	OtherSys      gauge
	PauseTotalNs  gauge
	StackInuse    gauge
	StackSys      gauge
	Sys           gauge
	TotalAlloc    gauge
	PollCount     counter
	RandomValue   gauge
}

type MetricsWrapper struct {
	metrics Metrics
	mu      sync.Mutex
}

func (metricsWrapper *MetricsWrapper) Set(metrics Metrics) {
	metricsWrapper.mu.Lock()
	defer metricsWrapper.mu.Unlock()
	metricsWrapper.metrics = metrics
}

func (metricsWrapper *MetricsWrapper) Get() Metrics {
	metricsWrapper.mu.Lock()
	defer metricsWrapper.mu.Unlock()
	return metricsWrapper.metrics
}
