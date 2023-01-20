package metrics

import (
	"math/rand"
	"runtime"

	"github.com/devldavydov/promytheus/internal/common/types"
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

func (c *RuntimeCollector) Collect() (Metrics, error) {
	c.collectCnt += 1

	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	metrics := Metrics{
		"Alloc":         types.Gauge(memStats.Alloc),
		"BuckHashSys":   types.Gauge(memStats.BuckHashSys),
		"Frees":         types.Gauge(memStats.Frees),
		"GCCPUFraction": types.Gauge(memStats.GCCPUFraction),
		"GCSys":         types.Gauge(memStats.GCSys),
		"HeapAlloc":     types.Gauge(memStats.HeapAlloc),
		"HeapIdle":      types.Gauge(memStats.HeapIdle),
		"HeapInuse":     types.Gauge(memStats.HeapInuse),
		"HeapObjects":   types.Gauge(memStats.HeapObjects),
		"HeapReleased":  types.Gauge(memStats.HeapReleased),
		"HeapSys":       types.Gauge(memStats.HeapSys),
		"LastGC":        types.Gauge(memStats.LastGC),
		"Lookups":       types.Gauge(memStats.Lookups),
		"MCacheInuse":   types.Gauge(memStats.MCacheInuse),
		"MCacheSys":     types.Gauge(memStats.MCacheSys),
		"MSpanInuse":    types.Gauge(memStats.MSpanInuse),
		"MSpanSys":      types.Gauge(memStats.MSpanSys),
		"Mallocs":       types.Gauge(memStats.Mallocs),
		"NextGC":        types.Gauge(memStats.NextGC),
		"NumForcedGC":   types.Gauge(memStats.NumForcedGC),
		"NumGC":         types.Gauge(memStats.NumGC),
		"OtherSys":      types.Gauge(memStats.OtherSys),
		"PauseTotalNs":  types.Gauge(memStats.PauseTotalNs),
		"StackInuse":    types.Gauge(memStats.StackInuse),
		"StackSys":      types.Gauge(memStats.StackSys),
		"Sys":           types.Gauge(memStats.Sys),
		"TotalAlloc":    types.Gauge(memStats.TotalAlloc),
		"PollCount":     types.Counter(c.collectCnt),
		"RandomValue":   types.Gauge(rand.Float64()),
	}

	c.logger.Debugf("Collected metrics: %+v", metrics)

	return metrics, nil
}
