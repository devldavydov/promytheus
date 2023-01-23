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
	Collect() (metric.Metrics, error)
}

type Publisher interface {
	Publish(metrics metric.Metrics) error
}

type Service struct {
	mu             sync.Mutex
	settings       ServiceSettings
	logger         *logrus.Logger
	currentMetrics metric.Metrics
	collector      Collector
	publisher      Publisher
}

func NewService(settings ServiceSettings, logger *logrus.Logger) *Service {
	return &Service{
		settings:  settings,
		logger:    logger,
		collector: collector.NewRuntimeCollector(logger),
		publisher: publisher.NewHTTPPublisher(settings.serverAddress, logger),
	}
}

func (service *Service) Start(ctx context.Context) error {
	service.logger.Info("Agent service started")

	var wg sync.WaitGroup

	wg.Add(2)
	go service.collectorThread(ctx, &wg)
	go service.publisherThread(ctx, &wg)
	wg.Wait()

	service.logger.Info("Agent service finished")
	return nil
}

func (service *Service) collectorThread(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(service.settings.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			service.logger.Debugf("Start collecting metrics")

			metrics, err := service.collector.Collect()
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

	ticker := time.NewTicker(service.settings.reportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			service.logger.Debugf("Start reporting metrics")

			service.mu.Lock()
			metrics := service.currentMetrics
			service.mu.Unlock()

			err := service.publisher.Publish(metrics)
			if err != nil {
				service.logger.Errorf("Publish metrics error: %s", err)
			}
		case <-ctx.Done():
			service.logger.Info("Publisher thread shutdown due to context closed")
			return
		}
	}
}
