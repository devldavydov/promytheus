// Package publisher is a package for different types of metric publishers.
package publisher

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/devldavydov/promytheus/internal/common/cipher"
	_http "github.com/devldavydov/promytheus/internal/common/http"
	"github.com/devldavydov/promytheus/internal/common/iotools"
	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/sirupsen/logrus"
)

const _httpClientTimeout = 1 * time.Second

// HTTPPublisher is a HTTP metric publisher.
type HTTPPublisher struct {
	serverAddress        *url.URL
	hmacKey              *string
	httpClient           *http.Client
	metricsChan          <-chan metric.Metrics
	logger               *logrus.Logger
	failedCounterMetrics metric.Metrics
	bufPool              *sync.Pool
	threadID             int
	shutdownTimeout      time.Duration
}

// NewHTTPPublisher creates new HTTPPublisher.
func NewHTTPPublisher(
	serverAddress *url.URL,
	hmacKey *string,
	metricsChan <-chan metric.Metrics,
	threadID int,
	cryptoPubKey *rsa.PublicKey,
	shutdownTimeout time.Duration,
	logger *logrus.Logger,
) *HTTPPublisher {
	client := &http.Client{
		Timeout: _httpClientTimeout,
	}

	bufPool := &sync.Pool{
		New: func() any {
			if cryptoPubKey == nil {
				return bytes.NewBuffer([]byte{})
			}
			return cipher.NewEncBuffer(cryptoPubKey)
		},
	}

	return &HTTPPublisher{
		serverAddress:   serverAddress,
		hmacKey:         hmacKey,
		httpClient:      client,
		metricsChan:     metricsChan,
		threadID:        threadID,
		bufPool:         bufPool,
		shutdownTimeout: shutdownTimeout,
		logger:          logger,
	}
}

func (httpPublisher *HTTPPublisher) Publish() {
	for metricsToSend := range httpPublisher.metricsChan {
		httpPublisher.processMetrics([]metric.Metrics{metricsToSend, httpPublisher.failedCounterMetrics})
	}
	// If channel closed, try to send failed metrics and exit
	httpPublisher.shutdown()
	httpPublisher.logger.Infof("Publisher[%d] thread shutdown due to context closed", httpPublisher.threadID)
	// for {
	// 	select {
	// 	case metricsToSend, ok := <-httpPublisher.metricsChan:
	// 		if !ok {
	// 			// If channel closed, try to send failed metrics and exit
	// 			httpPublisher.shutdown()
	// 			httpPublisher.logger.Infof("Publisher[%d] thread shutdown due to context closed", httpPublisher.threadID)
	// 			return
	// 		}
	// 		httpPublisher.processMetrics([]metric.Metrics{metricsToSend, httpPublisher.failedCounterMetrics})
	// 	}
	// }
}

func (httpPublisher *HTTPPublisher) processMetrics(metricsList []metric.Metrics) {
	var counterMetricsToSend = make(metric.Metrics)

	httpPublisher.logger.Debugf("Publisher[%d] publishing metrics: %+v", httpPublisher.threadID, metricsList)

	metricReq := make([]metric.MetricsDTO, 0, totalMetrics(metricsList))

	iterateMetrics(metricsList, func(name string, value metric.MetricValue) {
		metricReq = append(metricReq, httpPublisher.prepareMetric(name, value))

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
		httpPublisher.logger.Errorf("publisher[%d] failed to publish: %v", httpPublisher.threadID, err)
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

	request, err := http.NewRequest(
		http.MethodPost,
		httpPublisher.serverAddress.JoinPath("updates/").String(),
		buf)
	if err != nil {
		return fmt.Errorf("publisher[%d] failed to create request: %w", httpPublisher.threadID, err)
	}

	request.Header.Set("Content-Type", _http.ContentTypeApplicationJSON)

	response, err := httpPublisher.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("publisher[%d] failed to send request: %w", httpPublisher.threadID, err)
	}
	defer response.Body.Close()

	return nil
}

func (httpPublisher *HTTPPublisher) prepareMetric(metricName string, metricValue metric.MetricValue) metric.MetricsDTO {
	metricReq := metric.MetricsDTO{
		ID:    metricName,
		MType: metricValue.TypeName(),
	}

	if metric.GaugeTypeName == metricValue.TypeName() {
		metricReq.Value = metricValue.(metric.Gauge).FloatP()
	} else if metric.CounterTypeName == metricValue.TypeName() {
		metricReq.Delta = metricValue.(metric.Counter).IntP()
	}

	if httpPublisher.hmacKey != nil {
		hash := metricValue.Hmac(metricName, *httpPublisher.hmacKey)
		metricReq.Hash = &hash
	}

	return metricReq
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

func iterateMetrics(metricsList []metric.Metrics, fn func(name string, value metric.MetricValue)) {
	for _, metrics := range metricsList {
		for name, value := range metrics {
			fn(name, value)
		}
	}
}

func totalMetrics(metricsList []metric.Metrics) int {
	cnt := 0
	for _, m := range metricsList {
		cnt += len(m)
	}
	return cnt
}
