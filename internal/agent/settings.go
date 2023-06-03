package agent

import (
	"time"

	"github.com/devldavydov/promytheus/internal/common/nettools"
)

// ServiceSettings represents collecting metrics agent service settings.
type ServiceSettings struct {
	ServerAddress    nettools.Address
	HmacKey          *string
	CryptoPubKeyPath *string
	PollInterval     time.Duration
	ReportInterval   time.Duration
	RateLimit        int
	UseGRPC          bool
}

// NewServiceSettings creates new agent service settings.
func NewServiceSettings(
	serverAddress string,
	pollInterval time.Duration,
	reportInterval time.Duration,
	hmacKey string,
	rateLimit int,
	cryptoPubKeyPath string,
	useGRPC bool,
) (ServiceSettings, error) {
	srvAddr, err := nettools.NewAddress(serverAddress)
	if err != nil {
		return ServiceSettings{}, err
	}

	var hmac *string
	if hmacKey != "" {
		hmac = &hmacKey
	}

	var pubKeyPath *string
	if cryptoPubKeyPath != "" {
		pubKeyPath = &cryptoPubKeyPath
	}

	return ServiceSettings{
		ServerAddress:    srvAddr,
		PollInterval:     pollInterval,
		ReportInterval:   reportInterval,
		HmacKey:          hmac,
		RateLimit:        rateLimit,
		CryptoPubKeyPath: pubKeyPath,
		UseGRPC:          useGRPC,
	}, nil
}
