package server

import (
	"net"

	"github.com/devldavydov/promytheus/internal/common/nettools"
	"github.com/devldavydov/promytheus/internal/server/storage"
)

type ServiceSettings struct {
	HttpSettings      nettools.Address
	DatabaseDsn       string
	PersistSettings   storage.PersistSettings
	HmacKey           *string
	CryptoPrivKeyPath *string
	TrustedSubnet     *net.IPNet
	GrpcSettings      *nettools.Address
}

func NewServiceSettings(
	httpSettings nettools.Address,
	hmacKey string,
	databaseDsn string,
	persistSettimgs storage.PersistSettings,
	cryptoPrivKeyPath string,
	trustedSubnet *net.IPNet,
	grpcSettings *nettools.Address,
) ServiceSettings {
	var hmac *string
	if hmacKey != "" {
		hmac = &hmacKey
	}

	var privKeyPath *string
	if cryptoPrivKeyPath != "" {
		privKeyPath = &cryptoPrivKeyPath
	}

	return ServiceSettings{
		HttpSettings:      httpSettings,
		PersistSettings:   persistSettimgs,
		HmacKey:           hmac,
		DatabaseDsn:       databaseDsn,
		CryptoPrivKeyPath: privKeyPath,
		TrustedSubnet:     trustedSubnet,
		GrpcSettings:      grpcSettings,
	}
}
