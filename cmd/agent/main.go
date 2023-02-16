package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/devldavydov/promytheus/internal/agent"
	_log "github.com/devldavydov/promytheus/internal/common/log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config, err := LoadConfig(*flag.CommandLine, os.Args[1:])
	if err != nil {
		return fmt.Errorf("failed to load flag and ENV settings: %w", err)
	}

	logger, err := _log.NewLogger(config.LogLevel)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	agentSettings, err := AgentSettingsAdapt(config)
	if err != nil {
		return fmt.Errorf("failed to create agent settings: %w", err)
	}

	agentService := agent.NewService(agentSettings, logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	return agentService.Start(ctx)
}
