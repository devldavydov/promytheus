package agent

import (
	"net/url"
	"time"
)

type ServiceSettings struct {
	ServerAddress  *url.URL
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func NewServiceSettings(serverAddress string, pollInterval time.Duration, reportInterval time.Duration) (ServiceSettings, error) {
	url, err := url.ParseRequestURI(serverAddress)
	if err != nil {
		return ServiceSettings{}, err
	}

	return ServiceSettings{
		ServerAddress:  url,
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
	}, nil
}
