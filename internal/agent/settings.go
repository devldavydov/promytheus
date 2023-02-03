package agent

import (
	"net/url"
	"time"
)

type ServiceSettings struct {
	serverAddress  *url.URL
	pollInterval   time.Duration
	reportInterval time.Duration
}

func NewServiceSettings(serverAddress string, pollInterval time.Duration, reportInterval time.Duration) (ServiceSettings, error) {
	url, err := url.ParseRequestURI(serverAddress)
	if err != nil {
		return ServiceSettings{}, err
	}

	return ServiceSettings{
		serverAddress:  url,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
	}, nil
}

func (s ServiceSettings) GetServerAddress() *url.URL {
	return s.serverAddress
}

func (s ServiceSettings) GetPollInterval() time.Duration {
	return s.pollInterval
}

func (s ServiceSettings) GetReportInterval() time.Duration {
	return s.reportInterval
}
