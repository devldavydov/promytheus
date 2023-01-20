package agent

import (
	"context"
	"sync"
	"time"

	"github.com/devldavydov/promytheus/internal/agent/metrics"
	"github.com/sirupsen/logrus"
)

type Service struct {
	mu             sync.Mutex
	settings       ServiceSettings
	logger         *logrus.Logger
	currentMetrics metrics.Metrics
}

func NewService(settings ServiceSettings, logger *logrus.Logger) *Service {
	return &Service{settings: settings, logger: logger}
}

func (service *Service) Start(ctx context.Context) {
	service.logger.Info("Agent service started")

	var wg sync.WaitGroup

	wg.Add(2)
	go service.collectorThread(ctx, &wg)
	go service.publisherThread(ctx, &wg)
	wg.Wait()

	service.logger.Info("Agent service finished")
}

func (service *Service) collectorThread(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	collector := metrics.NewRuntimeCollector(service.logger)

	for {
		select {
		case <-time.After(service.settings.pollInterval):
			service.logger.Debugf("Start collecting metrics")

			metrics, err := collector.Collect()
			if err != nil {
				service.logger.Errorf("Failed to collect metrics, err: %s", err)
				continue
			}

			service.mu.Lock()
			service.currentMetrics = metrics
			service.mu.Unlock()
		case <-ctx.Done():
			service.logger.Info("Collector thread shutdown due to context closed")
			return
		}
	}
}

func (service *Service) publisherThread(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	publisher := metrics.NewHTTPPublisher(service.settings.serverAddress, service.logger)

	for {
		select {
		case <-time.After(service.settings.reportInterval):
			service.logger.Debugf("Start reporting metrics")

			service.mu.Lock()
			metrics := service.currentMetrics
			service.mu.Unlock()

			err := publisher.Publish(metrics)
			if err != nil {
				service.logger.Errorf("Publish metrics error: %s", err)
			}
		case <-ctx.Done():
			service.logger.Info("Publisher thread shutdown due to context closed")
			return
		}
	}
}
