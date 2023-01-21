package metrics

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

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
		Timeout: 1 * time.Second,
	}

	return &HTTPPublisher{serverAddress: serverAddress, httpClient: client, logger: logger}
}

func (httpPublisher *HTTPPublisher) Publish(metrics Metrics) error {
	metricsSentCnt := 0

	for name, value := range metrics {
		request, err := http.NewRequest(
			http.MethodPost,
			httpPublisher.serverAddress.JoinPath("update", value.TypeName(), name, value.String()).String(),
			nil)
		if err != nil {
			httpPublisher.logger.Errorf("Failed to create publish metrics request: %v", err)
			continue
		}
		request.Header.Set("Content-Type", "text/plain")

		response, err := httpPublisher.httpClient.Do(request)
		if err != nil {
			httpPublisher.logger.Errorf("Failed to publish metrics: %v", err)
			continue
		}
		response.Body.Close()

		metricsSentCnt += 1
	}

	if metricsSentCnt == 0 {
		return errors.New("failed to publish all metrics")
	}
	return nil
}
