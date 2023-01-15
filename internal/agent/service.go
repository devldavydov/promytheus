package agent

import (
	"context"
	"sync"
	"time"

	"github.com/devldavydov/promytheus/internal/agent/metrics"
	"github.com/sirupsen/logrus"
)

type Service struct {
	settings       ServiceSettings
	serviceCtx     context.Context
	logger         *logrus.Logger
	currentMetrics *metrics.MetricsWrapper
}

func NewService(ctx context.Context, settings ServiceSettings, logger *logrus.Logger) *Service {
	return &Service{serviceCtx: ctx, settings: settings, logger: logger, currentMetrics: &metrics.MetricsWrapper{}}
}

func (service *Service) Start() {
	service.logger.Info("Agent service started")

	var wg sync.WaitGroup

	wg.Add(2)
	go service.collectorThread(&wg, service.currentMetrics)
	go service.publisherThread(&wg, service.currentMetrics)
	wg.Wait()

	service.logger.Info("Agent service finished")
}

func (service *Service) collectorThread(wg *sync.WaitGroup, metricsWrapper *metrics.MetricsWrapper) {
	defer wg.Done()

	collector := metrics.NewRuntimeCollector(service.logger)

	for {
		select {
		case <-time.After(service.settings.pollInterval):
			service.logger.Debugf("Start collecting metrics")
			metricsVal, err := collector.Collect()
			if err != nil {
				service.logger.Errorf("Failed to collect metrics, err: %s", err)
			}
			metricsWrapper.Set(metricsVal)
		case <-service.serviceCtx.Done():
			service.logger.Info("Collector thread shutdown due to context closed")
			return
		}
	}
}

func (service *Service) publisherThread(wg *sync.WaitGroup, metricsWrapper *metrics.MetricsWrapper) {
	defer wg.Done()
	for {
		select {
		case <-time.After(service.settings.reportInterval):
			service.logger.Debugf("Start reporting metrics")
			service.logger.Debugf("Send metrics: %+v", metricsWrapper.Get())
		case <-service.serviceCtx.Done():
			service.logger.Info("Publisher thread shutdown due to context closed")
			return
		}
	}
}
