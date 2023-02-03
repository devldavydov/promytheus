package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devldavydov/promytheus/internal/common/logging"
	"github.com/devldavydov/promytheus/internal/server"
)

func main() {
	envConfig, err := LoadEnvConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load ENV settings: %v", err))
	}

	logger, err := logging.CreateLogger(envConfig.LogLevel)
	if err != nil {
		panic(err)
	}

	serverSettings, err := ServerSettingsAdapt(envConfig)
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
