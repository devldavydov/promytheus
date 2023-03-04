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
	Publish(context.Context)
}

type Service struct {
	settings                    ServiceSettings
	logger                      *logrus.Logger
	failedPublishCounterMetrics metric.Metrics
	collectors                  []Collector
	publisherFactory            func() Publisher
	metricsChan                 chan metric.Metrics
}

func NewService(settings ServiceSettings, logger *logrus.Logger) *Service {
	collectors := []Collector{
		collector.NewRuntimeCollector(settings.PollInterval, logger),
		collector.NewPsUtilCollector(settings.PollInterval, logger),
	}

	ch := make(chan metric.Metrics, len(collectors)*3)

	return &Service{
		settings:    settings,
		logger:      logger,
		collectors:  collectors,
		metricsChan: ch,
		publisherFactory: func() Publisher {
			return publisher.NewHTTPPublisher(settings.ServerAddress, settings.HmacKey, ch, logger)
		},
	}
}

func (service *Service) Start(ctx context.Context) error {
	service.logger.Info("Agent service started")

	var wg sync.WaitGroup

	for _, collector := range service.collectors {
		wg.Add(1)
		go func(ctx context.Context, clctr Collector) {
			defer wg.Done()
			clctr.Start(ctx)
		}(ctx, collector)
	}

	for i := 0; i < service.settings.RateLimit; i++ {
		wg.Add(1)
		go func(ctx context.Context) {
			defer wg.Done()
			service.publisherFactory().Publish(ctx)
		}(ctx)
	}

	service.startMainLoop(ctx)
	wg.Wait()

	service.logger.Info("Agent service finished")
	return nil
}

func (service *Service) startMainLoop(ctx context.Context) {
	ticker := time.NewTicker(service.settings.ReportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			service.logger.Debugf("Start reporting metrics")

			for _, collector := range service.collectors {
				metrics, err := collector.Collect()
				if err != nil {
					service.logger.Errorf("Failed to collect metrics from collector: %v", err)
					continue
				}

				service.metricsChan <- metrics
			}
		case <-ctx.Done():
			service.logger.Info("Main loop shutdown due to context closed")
			return
		}
	}
}
