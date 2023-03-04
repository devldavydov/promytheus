package collector

import (
	"context"
	"sync"
	"time"

	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/sirupsen/logrus"
)

type collectWorker interface {
	getMetrics() (metric.Metrics, error)
	collectCleanup()
}

type Collector struct {
	collectWorker
	mu             sync.Mutex
	pollInterval   time.Duration
	currentMetrics metric.Metrics
	logger         *logrus.Logger
}

func (c *Collector) Start(ctx context.Context) {
	ticker := time.NewTicker(c.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()

			metrics, err := c.getMetrics()
			if err != nil {
				c.logger.Errorf("Failed to get runtime metrics: %v", err)
				continue
			}
			c.currentMetrics = metrics

			c.mu.Unlock()
		case <-ctx.Done():
			c.logger.Info("Collector thread shutdown due to context closed")
			return
		}
	}
}

func (c *Collector) Collect() (metric.Metrics, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.collectCleanup()

	c.logger.Debugf("Collected metrics: %+v", c.currentMetrics)

	return c.currentMetrics, nil
}
