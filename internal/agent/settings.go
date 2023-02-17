package agent

import (
	"net/url"
	"time"
)

type ServiceSettings struct {
	ServerAddress  *url.URL
	PollInterval   time.Duration
	ReportInterval time.Duration
	HmacKey        *string
}

func NewServiceSettings(serverAddress string, pollInterval time.Duration, reportInterval time.Duration, hmacKey string) (ServiceSettings, error) {
	url, err := url.ParseRequestURI(serverAddress)
	if err != nil {
		return ServiceSettings{}, err
	}

	var hmac *string
	if hmacKey != "" {
		hmac = &hmacKey
	}

	return ServiceSettings{
		ServerAddress:  url,
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
		HmacKey:        hmac,
	}, nil
}
