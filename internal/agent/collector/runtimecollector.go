package collector

import (
	"context"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/sirupsen/logrus"
)

type RuntimeCollector struct {
	mu             sync.Mutex
	pollCnt        int64
	pollInterval   time.Duration
	currentMetrics metric.Metrics
	logger         *logrus.Logger
}

func NewRuntimeCollector(pollInterval time.Duration, logger *logrus.Logger) *RuntimeCollector {
	return &RuntimeCollector{pollCnt: 0, pollInterval: pollInterval, logger: logger}
}

func (rc *RuntimeCollector) Start(ctx context.Context) {
	ticker := time.NewTicker(rc.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rc.mu.Lock()

			metrics, err := rc.getRuntimeMetrics()
			if err != nil {
				rc.logger.Errorf("Failed to get runtime metrics: %v", err)
				continue
			}
			rc.currentMetrics = metrics

			rc.mu.Unlock()
		case <-ctx.Done():
			rc.logger.Info("Collector thread shutdown due to context closed")
			return
		}
	}
}

func (rc *RuntimeCollector) getRuntimeMetrics() (metric.Metrics, error) {
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

func (rc *RuntimeCollector) Collect() (metric.Metrics, error) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.pollCnt = 0

	rc.logger.Debugf("Collected metrics: %+v", rc.currentMetrics)

	return rc.currentMetrics, nil
}
