package main

import (
	"time"

	"github.com/devldavydov/promytheus/internal/agent"
)

func main() {
	agentSettings := agent.NewServiceSettings("127.0.0.1:8080", 2*time.Second, 10*time.Second)
	agentService := agent.NewService(agentSettings)
	agentService.Start()
}
