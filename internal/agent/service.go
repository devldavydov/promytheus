package agent

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Service struct {
	settings   ServiceSettings
	serviceCtx context.Context
	logger     *logrus.Logger
}

func NewService(ctx context.Context, settings ServiceSettings, logger *logrus.Logger) *Service {
	return &Service{serviceCtx: ctx, settings: settings, logger: logger}
}

func (service *Service) Start() {
	service.logger.Info("Agent service started")
	<-service.serviceCtx.Done()
	service.logger.Info("Agent service finished")
}
