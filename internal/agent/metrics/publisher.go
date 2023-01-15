package metrics

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const httpClientTimeout time.Duration = 1 * time.Second

type Publisher interface {
	Publish(metrics Metrics) error
}

type HttpPublisher struct {
	serverAddress string
	httpClient    *http.Client
	logger        *logrus.Logger
}

func NewHttpPublisher(serverAddress string, logger *logrus.Logger) *HttpPublisher {
	client := &http.Client{}
	client.Timeout = httpClientTimeout

	return &HttpPublisher{serverAddress: serverAddress, httpClient: client, logger: logger}
}

func (httpPublisher *HttpPublisher) Publish(metrics Metrics) error {
	urlFormat := fmt.Sprintf("http://%s/update/%%s/%%s/%%s", httpPublisher.serverAddress)
	metricsSentCnt := 0

	for _, m := range metrics.ToItemsList() {
		requestUrl := fmt.Sprintf(urlFormat, m.typeName, m.metricName, m.value)
		request, err := http.NewRequest(http.MethodPost, requestUrl, nil)
		if err != nil {
			httpPublisher.logger.Errorf("Failed to create publish metrics request: %v", err)
			continue
		}
		request.Header.Set("Content-Type", "text/plain")

		_, err = httpPublisher.httpClient.Do(request)
		if err != nil {
			httpPublisher.logger.Errorf("Failed to publish metrics: %v", err)
			continue
		}
		metricsSentCnt += 1
	}

	if metricsSentCnt == 0 {
		return errors.New("failed to publish all metrics")
	}
	return nil
}
