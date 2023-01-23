package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devldavydov/promytheus/internal/agent"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel) // TODO: read from env LOG_LEVEL

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	agentSettings, err := agent.NewServiceSettings("http://127.0.0.1:8080", 2*time.Second, 10*time.Second) // TODO: read settings from env/args
	if err != nil {
		logger.Errorf("Failed to create agent settings: %v", err)
		os.Exit(1)
	}
	agentService := agent.NewService(agentSettings, logger)

	err = agentService.Start(ctx)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
