package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/devldavydov/promytheus/internal/server"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel) // TODO: read from env LOG_LEVEL

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serverSettings := server.NewServiceSettings("127.0.0.1", 8080)
	serverService := server.NewService(serverSettings, 5*time.Second, logger)
	serverService.Start(ctx)
}
