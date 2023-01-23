package metrics

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/devldavydov/promytheus/internal/common/types"
	"github.com/sirupsen/logrus"
)

const httpClientTimeout = 1 * time.Second

type Publisher interface {
	Publish(metrics Metrics) error
}

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

func (httpPublisher *HTTPPublisher) Publish(metrics Metrics) error {
	metricsSentCnt := 0

	var err error
	for name, value := range metrics {
		err = httpPublisher.publishMetric(name, value)
		if err != nil {
			httpPublisher.logger.Errorf("Failed to publish metric: %v", err)
			continue
		}

		metricsSentCnt += 1
	}

	if metricsSentCnt == 0 {
		return errors.New("failed to publish all metrics")
	}
	return nil
}

func (httpPublisher *HTTPPublisher) publishMetric(metricName string, metricValue types.MetricValue) error {
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
