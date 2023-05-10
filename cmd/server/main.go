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

	"github.com/devldavydov/promytheus/internal/common/info"
	_log "github.com/devldavydov/promytheus/internal/common/log"
	"github.com/devldavydov/promytheus/internal/server"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	appVer := info.FormatVersion(buildVersion, buildDate, buildCommit)
	fmt.Println(appVer)

	config, err := LoadConfig(*flag.CommandLine, os.Args[1:])
	if err != nil {
		return fmt.Errorf("failed to load flag and ENV settings: %w", err)
	}

	logger, closer, err := _log.NewLogger(config.LogLevel, config.LogFile)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer closer.Close()

	serverSettings, err := ServerSettingsAdapt(config)
	if err != nil {
		return fmt.Errorf("failed to create server settings: %w", err)
	}

	logger.Info(appVer)
	serverService := server.NewService(serverSettings, 5*time.Second, logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	return serverService.Start(ctx)
}
