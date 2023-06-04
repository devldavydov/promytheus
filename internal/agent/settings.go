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
	GRPCCACertPath   *string
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
	grpcCACertPath string,
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

	var grpcCACert *string
	if grpcCACertPath != "" {
		grpcCACert = &grpcCACertPath
	}

	return ServiceSettings{
		ServerAddress:    srvAddr,
		PollInterval:     pollInterval,
		ReportInterval:   reportInterval,
		HmacKey:          hmac,
		RateLimit:        rateLimit,
		CryptoPubKeyPath: pubKeyPath,
		UseGRPC:          useGRPC,
		GRPCCACertPath:   grpcCACert,
	}, nil
}
