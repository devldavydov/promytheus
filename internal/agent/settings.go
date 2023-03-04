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
	RateLimit      int
}

func NewServiceSettings(serverAddress string, pollInterval time.Duration, reportInterval time.Duration, hmacKey string, rateLimit int) (ServiceSettings, error) {
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
		RateLimit:      rateLimit,
	}, nil
}
