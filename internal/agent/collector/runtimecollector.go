package collector

import (
	"math/rand"
	"runtime"

	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/sirupsen/logrus"
)

type RuntimeCollector struct {
	collectCnt int64
	logger     *logrus.Logger
}

func NewRuntimeCollector(logger *logrus.Logger) *RuntimeCollector {
	return &RuntimeCollector{collectCnt: 0, logger: logger}
}

func (c *RuntimeCollector) Collect() (metric.Metrics, error) {
	c.collectCnt += 1

	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	metrics := metric.Metrics{
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
		"PollCount":     metric.Counter(c.collectCnt),
		"RandomValue":   metric.Gauge(rand.Float64()),
	}

	c.logger.Debugf("Collected metrics: %+v", metrics)

	return metrics, nil
}
