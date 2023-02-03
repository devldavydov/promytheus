package server

type ServiceSettings struct {
	serverAddress string
	serverPort    int
}

func NewServiceSettings(serverAddress string, serverPort int) ServiceSettings {
	return ServiceSettings{serverAddress: serverAddress, serverPort: serverPort}
}

func (s ServiceSettings) GetServerAddress() string {
	return s.serverAddress
}

func (s ServiceSettings) GetServerPort() int {
	return s.serverPort
}
