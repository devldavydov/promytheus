package agent

import (
	"net/url"
	"time"
)

// ServiceSettings represents collecting metrics agent service settings.
type ServiceSettings struct {
	ServerAddress  *url.URL
	HmacKey        *string
	PollInterval   time.Duration
	ReportInterval time.Duration
	RateLimit      int
}

// NewServiceSettings creates new agent service settings.
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
