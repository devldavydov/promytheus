package main

import (
	"context"
	"os"
	"time"

	"github.com/devldavydov/promytheus/internal/agent"
	"github.com/devldavydov/promytheus/internal/common/helpers"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel) // TODO: read from env LOG_LEVEL

	agentSettings, err := agent.NewServiceSettings("http://127.0.0.1:8080", 2*time.Second, 10*time.Second) // TODO: read settings from env/args
	if err != nil {
		logger.Errorf("Failed to create agent settings: %v", err)
		os.Exit(1)
	}

	agentService := agent.NewService(helpers.CreateContextWithSignalHadler(context.Background()), agentSettings, logger)
	agentService.Start()
}
