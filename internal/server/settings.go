package server

import "github.com/devldavydov/promytheus/internal/server/storage"

type ServiceSettings struct {
	ServerAddress   string
	ServerPort      int
	PersistSettings storage.PersistSettings
}

func NewServiceSettings(serverAddress string, serverPort int, persistSettimgs storage.PersistSettings) ServiceSettings {
	return ServiceSettings{
		ServerAddress:   serverAddress,
		ServerPort:      serverPort,
		PersistSettings: persistSettimgs}
}
