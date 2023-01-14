package agent

import "time"

type ServiceSettings struct {
	serverAddress  string
	pollInterval   time.Duration
	reportInterval time.Duration
}

func NewServiceSettings(serverAddress string, pollInterval time.Duration, reportInterval time.Duration) ServiceSettings {
	return ServiceSettings{
		serverAddress:  serverAddress,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
	}
}
