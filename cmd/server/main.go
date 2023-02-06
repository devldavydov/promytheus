package main

import (
	"context"
	"flag"
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

	flagConfig, err := LoadFlagConfig(*flag.CommandLine, os.Args[1:])
	if err != nil {
		panic(fmt.Sprintf("Failed to load flag settings: %v", err))
	}

	logger, err := logging.CreateLogger(envConfig.LogLevel.Value)
	if err != nil {
		panic(err)
	}

	serverSettings, err := ServerSettingsAdapt(envConfig, flagConfig)
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
