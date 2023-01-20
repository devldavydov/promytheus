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
