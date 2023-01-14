package agent

import "fmt"

type Service struct {
	settings ServiceSettings
}

func NewService(settings ServiceSettings) *Service {
	return &Service{settings: settings}
}

func (service *Service) Start() {
	fmt.Println("Agent started")
}
