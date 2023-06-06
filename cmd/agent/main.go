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

	"github.com/devldavydov/promytheus/internal/agent"
	"github.com/devldavydov/promytheus/internal/common/info"
	_log "github.com/devldavydov/promytheus/internal/common/log"
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
		return fmt.Errorf("failed to load configuration settings: %w", err)
	}

	logger, closer, err := _log.NewLogger(config.LogLevel, config.LogFile)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer closer.Close()

	agentSettings, err := AgentSettingsAdapt(config)
	if err != nil {
		return fmt.Errorf("failed to create agent settings: %w", err)
	}

	logger.Info(appVer)
	agentService, err := agent.NewService(agentSettings, 5*time.Second, logger)
	if err != nil {
		return fmt.Errorf("failed to create agent service: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	return agentService.Start(ctx)
}
