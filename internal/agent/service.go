// Package agent is the main package for agent service.
package agent

import (
	"context"
	"crypto/rsa"
	"sync"
	"time"

	"github.com/devldavydov/promytheus/internal/agent/collector"
	"github.com/devldavydov/promytheus/internal/agent/publisher"
	"github.com/devldavydov/promytheus/internal/common/cipher"
	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/devldavydov/promytheus/internal/common/nettools"
	"github.com/devldavydov/promytheus/internal/grpc/gtls"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials"
)

// Collector is an interface for collector functionality.
type Collector interface {
	Start(context.Context)
	Collect() (metric.Metrics, error)
}

// Publisher is an interface for publisher functionality.
type Publisher interface {
	Publish()
}

// Service represents collecting metrics agent service.
type Service struct {
	logger                      *logrus.Logger
	failedPublishCounterMetrics metric.Metrics
	publisherFactory            PublisherFactory
	metricsChan                 chan metric.Metrics
	collectors                  []Collector
	settings                    ServiceSettings
}

// NewService creates new agent service.
func NewService(settings ServiceSettings, shutdownTimeout time.Duration, logger *logrus.Logger) (*Service, error) {
	collectors := []Collector{
		collector.NewRuntimeCollector(settings.PollInterval, logger),
		collector.NewPsUtilCollector(settings.PollInterval, logger),
	}

	ch := make(chan metric.Metrics, len(collectors)*2)

	hostIP, err := nettools.GetHostIP()
	if err != nil {
		return nil, err
	}

	return &Service{
		settings:    settings,
		logger:      logger,
		collectors:  collectors,
		metricsChan: ch,
		publisherFactory: CreatePublisherFactory(
			settings, shutdownTimeout, ch, hostIP, logger,
		),
	}, nil
}

// Start runs agent service with context.
func (service *Service) Start(ctx context.Context) error {
	service.logger.Info("Agent service started")

	// Load encryption settings
	encrSettings, err := service.loadEncryptionSettings()
	if err != nil {
		return err
	}

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
		go func(ctx context.Context, threadID int) {
			defer wg.Done()
			service.publisherFactory(threadID, encrSettings).Publish()
		}(ctx, i+1)
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
			close(service.metricsChan)
			service.logger.Info("Main loop shutdown due to context closed")
			return
		}
	}
}

func (service *Service) loadEncryptionSettings() (publisher.EncryptionSettings, error) {
	var err error
	encrSettings := publisher.EncryptionSettings{}

	if !service.settings.UseGRPC {
		var cryptoPubKey *rsa.PublicKey
		cryptoPubKey, err = service.loadHTTPCryptoPubKey()
		if err == nil {
			encrSettings.CryptoPubKey = cryptoPubKey
		}
	} else {
		var tlsCredentials credentials.TransportCredentials
		tlsCredentials, err = service.loadGRPCTLS()
		if err == nil {
			encrSettings.TLSCredentials = tlsCredentials
		}
	}

	return encrSettings, err
}

func (service *Service) loadHTTPCryptoPubKey() (*rsa.PublicKey, error) {
	if service.settings.CryptoPubKeyPath == nil {
		return nil, nil
	}
	return cipher.PublicKeyFromFile(*service.settings.CryptoPubKeyPath)
}

func (service *Service) loadGRPCTLS() (credentials.TransportCredentials, error) {
	if service.settings.GRPCCACertPath == nil {
		return nil, nil
	}
	return gtls.LoadCACert(*service.settings.GRPCCACertPath, "")
}
