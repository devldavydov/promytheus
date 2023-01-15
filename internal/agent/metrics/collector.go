package metrics

import (
	"math/rand"
	"runtime"

	"github.com/sirupsen/logrus"
)

type Collector interface {
	Collect() (Metrics, error)
}

type RuntimeCollector struct {
	collectCnt int64
	logger     *logrus.Logger
}

func NewRuntimeCollector(logger *logrus.Logger) *RuntimeCollector {
	return &RuntimeCollector{collectCnt: 0, logger: logger}
}

func (runtimeColl *RuntimeCollector) Collect() (Metrics, error) {
	runtimeColl.collectCnt += 1

	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	metrics := Metrics{
		Alloc:         gauge(memStats.Alloc),
		BuckHashSys:   gauge(memStats.BuckHashSys),
		Frees:         gauge(memStats.Frees),
		GCCPUFraction: gauge(memStats.GCCPUFraction),
		GCSys:         gauge(memStats.GCSys),
		HeapAlloc:     gauge(memStats.HeapAlloc),
		HeapIdle:      gauge(memStats.HeapIdle),
		HeapInuse:     gauge(memStats.HeapInuse),
		HeapObjects:   gauge(memStats.HeapObjects),
		HeapReleased:  gauge(memStats.HeapReleased),
		HeapSys:       gauge(memStats.HeapSys),
		LastGC:        gauge(memStats.LastGC),
		Lookups:       gauge(memStats.Lookups),
		MCacheInuse:   gauge(memStats.MCacheInuse),
		MCacheSys:     gauge(memStats.MCacheSys),
		MSpanInuse:    gauge(memStats.MSpanInuse),
		MSpanSys:      gauge(memStats.MSpanSys),
		Mallocs:       gauge(memStats.Mallocs),
		NextGC:        gauge(memStats.NextGC),
		NumForcedGC:   gauge(memStats.NumForcedGC),
		NumGC:         gauge(memStats.NumGC),
		OtherSys:      gauge(memStats.OtherSys),
		PauseTotalNs:  gauge(memStats.PauseTotalNs),
		StackInuse:    gauge(memStats.StackInuse),
		StackSys:      gauge(memStats.StackSys),
		Sys:           gauge(memStats.Sys),
		TotalAlloc:    gauge(memStats.TotalAlloc),
		PollCount:     counter(runtimeColl.collectCnt),
		RandomValue:   gauge(rand.Float64()),
	}

	runtimeColl.logger.Debugf("Collected metrics: %+v", metrics)

	return metrics, nil
}
