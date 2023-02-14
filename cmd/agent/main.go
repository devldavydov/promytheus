package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/devldavydov/promytheus/internal/agent"
	"github.com/devldavydov/promytheus/internal/common/logging"
)

func main() {
	config, err := LoadConfig(*flag.CommandLine, os.Args[1:])
	if err != nil {
		panic(fmt.Sprintf("Failed to load flag and ENV settings: %v", err))
	}

	logger, err := logging.CreateLogger(config.LogLevel)
	if err != nil {
		panic(err)
	}

	agentSettings, err := AgentSettingsAdapt(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create agent settings: %v", err))
	}

	agentService := agent.NewService(agentSettings, logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err = agentService.Start(ctx)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
