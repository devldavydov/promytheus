package collector

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/sirupsen/logrus"
)

type RuntimeCollector struct {
	pollCnt int64
}

var _ collectWorker = (*RuntimeCollector)(nil)

func NewRuntimeCollector(pollInterval time.Duration, logger *logrus.Logger) *Collector {
	return &Collector{
		collectWorker: &RuntimeCollector{pollCnt: 0},
		name:          "RuntimeCollector",
		pollInterval:  pollInterval,
		logger:        logger,
	}
}

func (rc *RuntimeCollector) getMetrics() (metric.Metrics, error) {
	rc.pollCnt += 1
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	return metric.Metrics{
		"Alloc":         metric.Gauge(memStats.Alloc),
		"BuckHashSys":   metric.Gauge(memStats.BuckHashSys),
		"Frees":         metric.Gauge(memStats.Frees),
		"GCCPUFraction": metric.Gauge(memStats.GCCPUFraction),
		"GCSys":         metric.Gauge(memStats.GCSys),
		"HeapAlloc":     metric.Gauge(memStats.HeapAlloc),
		"HeapIdle":      metric.Gauge(memStats.HeapIdle),
		"HeapInuse":     metric.Gauge(memStats.HeapInuse),
		"HeapObjects":   metric.Gauge(memStats.HeapObjects),
		"HeapReleased":  metric.Gauge(memStats.HeapReleased),
		"HeapSys":       metric.Gauge(memStats.HeapSys),
		"LastGC":        metric.Gauge(memStats.LastGC),
		"Lookups":       metric.Gauge(memStats.Lookups),
		"MCacheInuse":   metric.Gauge(memStats.MCacheInuse),
		"MCacheSys":     metric.Gauge(memStats.MCacheSys),
		"MSpanInuse":    metric.Gauge(memStats.MSpanInuse),
		"MSpanSys":      metric.Gauge(memStats.MSpanSys),
		"Mallocs":       metric.Gauge(memStats.Mallocs),
		"NextGC":        metric.Gauge(memStats.NextGC),
		"NumForcedGC":   metric.Gauge(memStats.NumForcedGC),
		"NumGC":         metric.Gauge(memStats.NumGC),
		"OtherSys":      metric.Gauge(memStats.OtherSys),
		"PauseTotalNs":  metric.Gauge(memStats.PauseTotalNs),
		"StackInuse":    metric.Gauge(memStats.StackInuse),
		"StackSys":      metric.Gauge(memStats.StackSys),
		"Sys":           metric.Gauge(memStats.Sys),
		"TotalAlloc":    metric.Gauge(memStats.TotalAlloc),
		"PollCount":     metric.Counter(rc.pollCnt),
		"RandomValue":   metric.Gauge(rand.Float64()),
	}, nil
}

func (rc *RuntimeCollector) collectCleanup() {
	rc.pollCnt = 0
}
