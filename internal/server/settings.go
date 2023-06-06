package server

import (
	"net"

	"github.com/devldavydov/promytheus/internal/common/nettools"
	"github.com/devldavydov/promytheus/internal/grpc/gtls"
	"github.com/devldavydov/promytheus/internal/server/storage"
)

type ServiceSettings struct {
	HTTPAddress       nettools.Address
	DatabaseDsn       string
	PersistSettings   storage.PersistSettings
	HmacKey           *string
	CryptoPrivKeyPath *string
	TrustedSubnet     *net.IPNet
	GRPCAddress       *nettools.Address
	GRPCServerTLS     *gtls.TLSServerSettings
}

func NewServiceSettings(
	httpAddress nettools.Address,
	hmacKey string,
	databaseDsn string,
	persistSettimgs storage.PersistSettings,
	cryptoPrivKeyPath string,
	trustedSubnet *net.IPNet,
	grpcAddress *nettools.Address,
	grpcServerTLS *gtls.TLSServerSettings,
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
		HTTPAddress:       httpAddress,
		PersistSettings:   persistSettimgs,
		HmacKey:           hmac,
		DatabaseDsn:       databaseDsn,
		CryptoPrivKeyPath: privKeyPath,
		TrustedSubnet:     trustedSubnet,
		GRPCAddress:       grpcAddress,
		GRPCServerTLS:     grpcServerTLS,
	}
}
