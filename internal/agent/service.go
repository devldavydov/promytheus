package agent

import (
	"context"
	"sync"
	"time"

	"github.com/devldavydov/promytheus/internal/agent/collector"
	"github.com/devldavydov/promytheus/internal/agent/publisher"
	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/sirupsen/logrus"
)

type Collector interface {
	Start(context.Context)
	Collect() (metric.Metrics, error)
}

type Publisher interface {
	Publish([]metric.Metrics) (metric.Metrics, error)
}

type Service struct {
	settings                    ServiceSettings
	logger                      *logrus.Logger
	failedPublishCounterMetrics metric.Metrics
	collector                   Collector
	publisher                   Publisher
}

func NewService(settings ServiceSettings, logger *logrus.Logger) *Service {
	return &Service{
		settings:  settings,
		logger:    logger,
		collector: collector.NewRuntimeCollector(settings.pollInterval, logger),
		publisher: publisher.NewHTTPPublisher(settings.serverAddress, logger),
	}
}

func (service *Service) Start(ctx context.Context) error {
	service.logger.Info("Agent service started")

	var wg sync.WaitGroup

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		service.collector.Start(ctx)
	}(ctx)

	service.startMainLoop(ctx)
	wg.Wait()

	service.logger.Info("Agent service finished")
	return nil
}

func (service *Service) startMainLoop(ctx context.Context) {
	ticker := time.NewTicker(service.settings.reportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			service.logger.Debugf("Start reporting metrics")

			metrics, err := service.collector.Collect()
			if err != nil {
				service.logger.Errorf("Failed to collect metrics from collector: %v", err)
				continue
			}

			failedPublishCounterMetrics, err := service.publisher.Publish([]metric.Metrics{metrics, service.failedPublishCounterMetrics})
			if err != nil {
				service.logger.Errorf("Publish metrics error: %v", err)
				service.failedPublishCounterMetrics = failedPublishCounterMetrics
			} else {
				service.failedPublishCounterMetrics = nil
			}
		case <-ctx.Done():
			service.logger.Info("Main loop shutdown due to context closed")
			return
		}
	}
}
