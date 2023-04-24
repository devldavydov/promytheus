// Package collector is a package for different types collectors.
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

// Collector is a base struct for different collectors.
type Collector struct {
	collectWorker
	mu             sync.Mutex
	name           string
	pollInterval   time.Duration
	currentMetrics metric.Metrics
	logger         *logrus.Logger
}

// Start - runs collector.
func (c *Collector) Start(ctx context.Context) {
	ticker := time.NewTicker(c.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()

			metrics, err := c.getMetrics()
			if err != nil {
				c.logger.Errorf("Collector [%s] failed to get runtime metrics: %v", c.name, err)
				continue
			}
			c.currentMetrics = metrics

			c.mu.Unlock()
		case <-ctx.Done():
			c.logger.Infof("Collector [%s] thread shutdown due to context closed", c.name)
			return
		}
	}
}

// Collect - collects metrics.
func (c *Collector) Collect() (metric.Metrics, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.collectCleanup()

	c.logger.Debugf("Collector [%s] collected metrics: %+v", c.name, c.currentMetrics)

	return c.currentMetrics, nil
}
