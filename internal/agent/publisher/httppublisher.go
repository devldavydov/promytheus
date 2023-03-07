package publisher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	_http "github.com/devldavydov/promytheus/internal/common/http"
	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/sirupsen/logrus"
)

const _httpClientTimeout = 1 * time.Second

type HTTPPublisher struct {
	serverAddress        *url.URL
	hmacKey              *string
	httpClient           *http.Client
	metricsChan          <-chan metric.Metrics
	threadID             int
	logger               *logrus.Logger
	failedCounterMetrics metric.Metrics
}

func NewHTTPPublisher(serverAddress *url.URL, hmacKey *string, metricsChan <-chan metric.Metrics, threadID int, logger *logrus.Logger) *HTTPPublisher {
	client := &http.Client{
		Timeout: _httpClientTimeout,
	}

	return &HTTPPublisher{
		serverAddress: serverAddress,
		hmacKey:       hmacKey,
		httpClient:    client,
		metricsChan:   metricsChan,
		threadID:      threadID,
		logger:        logger,
	}
}

func (httpPublisher *HTTPPublisher) Publish(ctx context.Context) {
	for {
		select {
		case metricsToSend := <-httpPublisher.metricsChan:
			httpPublisher.processMetrics(ctx, []metric.Metrics{metricsToSend, httpPublisher.failedCounterMetrics})
		case <-ctx.Done():
			httpPublisher.logger.Infof("Publisher[%d] thread shutdown due to context closed", httpPublisher.threadID)
			return
		}
	}
}

func (httpPublisher *HTTPPublisher) processMetrics(ctx context.Context, metricsList []metric.Metrics) {
	var counterMetricsToSend = make(metric.Metrics)

	httpPublisher.logger.Debugf("Publisher[%d] publishing metrics: %+v", httpPublisher.threadID, metricsList)

	metricReq := make([]metric.MetricsDTO, 0, totalMetrics(metricsList))

	iterateMetrics(ctx, metricsList, func(name string, value metric.MetricValue) {
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
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(metricReq)

	request, err := http.NewRequest(
		http.MethodPost,
		httpPublisher.serverAddress.JoinPath("updates/").String(),
		&buf)
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

func iterateMetrics(ctx context.Context, metricsList []metric.Metrics, fn func(name string, value metric.MetricValue)) {
	for _, metrics := range metricsList {
		for name, value := range metrics {
			if ctx.Err() != nil {
				return
			}
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
