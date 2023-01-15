package main

import (
	"context"
	"time"

	"github.com/devldavydov/promytheus/internal/agent"
	"github.com/devldavydov/promytheus/internal/common/helpers"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel) // TODO: read from env LOG_LEVEL

	agentSettings := agent.NewServiceSettings("127.0.0.1:8080", 2*time.Second, 10*time.Second)
	agentService := agent.NewService(helpers.CreateContextWithSignalHadler(context.Background(), logger), agentSettings, logger)
	agentService.Start()
}
