package publisher

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/sirupsen/logrus"
)

const httpClientTimeout = 1 * time.Second

type HTTPPublisher struct {
	serverAddress *url.URL
	httpClient    *http.Client
	logger        *logrus.Logger
}

func NewHTTPPublisher(serverAddress *url.URL, logger *logrus.Logger) *HTTPPublisher {
	client := &http.Client{
		Timeout: httpClientTimeout,
	}

	return &HTTPPublisher{serverAddress: serverAddress, httpClient: client, logger: logger}
}

func (httpPublisher *HTTPPublisher) Publish(metrics metric.Metrics) (metric.Metrics, error) {
	var failedPublishMetrics = make(metric.Metrics)
	var err error

	httpPublisher.logger.Debugf("Publishing metrics: %+v", metrics)

	for name, value := range metrics {
		err = httpPublisher.publishMetric(name, value)
		if err != nil {
			failedPublishMetrics[name] = value
		}
	}

	if len(failedPublishMetrics) != 0 {
		err = fmt.Errorf("failed to publish: %+v", failedPublishMetrics)
	}
	return failedPublishMetrics, err
}

func (httpPublisher *HTTPPublisher) publishMetric(metricName string, metricValue metric.MetricValue) error {
	request, err := http.NewRequest(
		http.MethodPost,
		httpPublisher.serverAddress.JoinPath("update", metricValue.TypeName(), metricName, metricValue.String()).String(),
		nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	request.Header.Set("Content-Type", "text/plain")

	response, err := httpPublisher.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	return nil
}
