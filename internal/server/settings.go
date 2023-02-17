package server

import "github.com/devldavydov/promytheus/internal/server/storage"

type ServiceSettings struct {
	ServerAddress   string
	ServerPort      int
	HmacKey         *string
	PersistSettings storage.PersistSettings
}

func NewServiceSettings(serverAddress string, serverPort int, hmacKey string, persistSettimgs storage.PersistSettings) ServiceSettings {
	var hmac *string
	if hmacKey != "" {
		hmac = &hmacKey
	}

	return ServiceSettings{
		ServerAddress:   serverAddress,
		ServerPort:      serverPort,
		PersistSettings: persistSettimgs,
		HmacKey:         hmac,
	}
}
