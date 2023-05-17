package server

import "github.com/devldavydov/promytheus/internal/server/storage"

type ServiceSettings struct {
	HmacKey           *string
	ServerAddress     string
	DatabaseDsn       string
	PersistSettings   storage.PersistSettings
	ServerPort        int
	CryptoPrivKeyPath *string
}

func NewServiceSettings(
	serverAddress string,
	serverPort int,
	hmacKey string,
	databaseDsn string,
	persistSettimgs storage.PersistSettings,
	cryptoPrivKeyPath string,
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
		ServerAddress:     serverAddress,
		ServerPort:        serverPort,
		PersistSettings:   persistSettimgs,
		HmacKey:           hmac,
		DatabaseDsn:       databaseDsn,
		CryptoPrivKeyPath: privKeyPath,
	}
}
