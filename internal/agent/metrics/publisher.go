package metrics

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	httpClientTimeout time.Duration = 1 * time.Second
	updateURLFormat   string        = "http://%s/update/%%s/%%s/%%s"
)

type Publisher interface {
	Publish(metrics Metrics) error
}

type HTTPPublisher struct {
	serverAddress string
	httpClient    *http.Client
	logger        *logrus.Logger
}

func NewHTTPPublisher(serverAddress string, logger *logrus.Logger) *HTTPPublisher {
	client := &http.Client{}
	client.Timeout = httpClientTimeout

	return &HTTPPublisher{serverAddress: serverAddress, httpClient: client, logger: logger}
}

func (httpPublisher *HTTPPublisher) Publish(metrics Metrics) error {
	urlFormat := fmt.Sprintf(updateURLFormat, httpPublisher.serverAddress)
	metricsSentCnt := 0

	for _, m := range metrics.ToItemsList() {
		requestURL := fmt.Sprintf(urlFormat, m.typeName, m.metricName, m.value)
		request, err := http.NewRequest(http.MethodPost, requestURL, nil)
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
		defer response.Body.Close()

		metricsSentCnt += 1
	}

	if metricsSentCnt == 0 {
		return errors.New("failed to publish all metrics")
	}
	return nil
}
