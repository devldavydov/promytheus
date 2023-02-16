package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_log "github.com/devldavydov/promytheus/internal/common/log"
	"github.com/devldavydov/promytheus/internal/server"
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

	serverSettings, err := ServerSettingsAdapt(config)
	if err != nil {
		return fmt.Errorf("failed to create server settings: %w", err)
	}

	serverService := server.NewService(serverSettings, 5*time.Second, logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	return serverService.Start(ctx)
}
