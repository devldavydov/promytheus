package server

import "github.com/devldavydov/promytheus/internal/server/storage"

type ServiceSettings struct {
	serverAddress   string
	serverPort      int
	persistSettings storage.PersistSettings
}

func NewServiceSettings(serverAddress string, serverPort int, persistSettimgs storage.PersistSettings) ServiceSettings {
	return ServiceSettings{
		serverAddress:   serverAddress,
		serverPort:      serverPort,
		persistSettings: persistSettimgs}
}

func (s ServiceSettings) GetServerAddress() string {
	return s.serverAddress
}

func (s ServiceSettings) GetServerPort() int {
	return s.serverPort
}

func (s ServiceSettings) GetPersistenSettings() storage.PersistSettings {
	return s.persistSettings
}
