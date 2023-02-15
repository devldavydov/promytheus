package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devldavydov/promytheus/internal/common/log"
	"github.com/devldavydov/promytheus/internal/server"
)

func main() {
	config, err := LoadConfig(*flag.CommandLine, os.Args[1:])
	if err != nil {
		panic(fmt.Sprintf("Failed to load flag and ENV settings: %v", err))
	}

	logger, err := log.NewLogger(config.LogLevel)
	if err != nil {
		panic(err)
	}

	serverSettings, err := ServerSettingsAdapt(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create server settings: %v", err))
	}

	serverService := server.NewService(serverSettings, 5*time.Second, logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err = serverService.Start(ctx)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
