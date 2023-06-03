// Package publisher is a package for different types of metric publishers.
package publisher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/devldavydov/promytheus/internal/common/cipher"
	_http "github.com/devldavydov/promytheus/internal/common/http"
	"github.com/devldavydov/promytheus/internal/common/iotools"
	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/devldavydov/promytheus/internal/common/nettools"
	"github.com/sirupsen/logrus"
)

const (
	_httpClientTimeout = 1 * time.Second
)

// HTTPPublisher is a HTTP metric publisher.
type HTTPPublisher struct {
	serverAddress        nettools.Address
	hmacKey              *string
	httpClient           *http.Client
	metricsChan          <-chan metric.Metrics
	logger               *logrus.Logger
	failedCounterMetrics metric.Metrics
	bufPool              *sync.Pool
	hostIP               string
	threadID             int
	shutdownTimeout      time.Duration
}

// NewHTTPPublisher creates new HTTPPublisher.
func NewHTTPPublisher(
	serverAddress nettools.Address,
	metricsChan <-chan metric.Metrics,
	threadID int,
	logger *logrus.Logger,
	extra PublisherExtraSettings,
) *HTTPPublisher {
	client := &http.Client{
		Timeout: _httpClientTimeout,
	}

	bufPool := &sync.Pool{
		New: func() any {
			if extra.EncrSettings.CryptoPubKey == nil {
				return bytes.NewBuffer([]byte{})
			}
			return cipher.NewEncBuffer(extra.EncrSettings.CryptoPubKey)
		},
	}

	shutdownTimeout := _defaultShutdownTimeout
	if extra.ShutdownTimeout != nil {
		shutdownTimeout = *extra.ShutdownTimeout
	}

	return &HTTPPublisher{
		serverAddress:   serverAddress,
		hmacKey:         extra.HmacKey,
		httpClient:      client,
		metricsChan:     metricsChan,
		threadID:        threadID,
		bufPool:         bufPool,
		shutdownTimeout: shutdownTimeout,
		hostIP:          extra.HostIP.String(),
		logger:          logger,
	}
}

func (httpPublisher *HTTPPublisher) Publish() {
	for metricsToSend := range httpPublisher.metricsChan {
		httpPublisher.processMetrics([]metric.Metrics{metricsToSend, httpPublisher.failedCounterMetrics})
	}
	// If channel closed, try to send failed metrics and exit
	httpPublisher.shutdown()
	httpPublisher.logger.Infof("HTTP publisher[%d] thread shutdown due to context closed", httpPublisher.threadID)
}

func (httpPublisher *HTTPPublisher) processMetrics(metricsList []metric.Metrics) {
	var counterMetricsToSend = make(metric.Metrics)

	httpPublisher.logger.Debugf("HTTP publisher[%d] publishing metrics: %+v", httpPublisher.threadID, metricsList)

	metricReq := make([]metric.MetricsDTO, 0, totalMetrics(metricsList))

	iterateMetrics(metricsList, func(name string, value metric.MetricValue) {
		metricReq = append(metricReq, prepareMetric(name, value, httpPublisher.hmacKey))

		if value.TypeName() == metric.CounterTypeName {
			curVal, ok := counterMetricsToSend[name]
			if !ok {
				counterMetricsToSend[name] = value
			} else {
				counterMetricsToSend[name] = curVal.(metric.Counter) + value.(metric.Counter)
			}
		}
	})

	if err := httpPublisher.publishMetrics(metricReq); err != nil {
		httpPublisher.logger.Errorf("HTTP publisher[%d] failed to publish: %v", httpPublisher.threadID, err)
		httpPublisher.failedCounterMetrics = counterMetricsToSend
		return
	}

	httpPublisher.failedCounterMetrics = nil
}

func (httpPublisher *HTTPPublisher) publishMetrics(metricReq []metric.MetricsDTO) error {
	buf := httpPublisher.bufPool.Get().(iotools.PoolBuffer)
	defer httpPublisher.bufPool.Put(buf)

	buf.Reset()
	json.NewEncoder(buf).Encode(metricReq)

	ctx, cancel := context.WithTimeout(context.Background(), _defaultRequestTimeout)
	defer cancel()

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("http://%s/updates/", httpPublisher.serverAddress.String()),
		buf)
	if err != nil {
		return fmt.Errorf("HTTP publisher[%d] failed to create request: %w", httpPublisher.threadID, err)
	}

	request.Header.Set("Content-Type", _http.ContentTypeApplicationJSON)
	request.Header.Set(nettools.RealIPHeader, httpPublisher.hostIP)

	response, err := httpPublisher.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("publisher[%d] failed to send request, err: %w", httpPublisher.threadID, err)
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("publisher[%d] failed to send request, code: %d", httpPublisher.threadID, response.StatusCode)
	}
	defer response.Body.Close()

	return nil
}

func (httpPublisher *HTTPPublisher) shutdown() {
	if httpPublisher.failedCounterMetrics == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), httpPublisher.shutdownTimeout)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if httpPublisher.failedCounterMetrics == nil {
				return
			}

			httpPublisher.processMetrics([]metric.Metrics{httpPublisher.failedCounterMetrics})
		case <-ctx.Done():
			return
		}
	}
}
