package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/devldavydov/promytheus/internal/server/handlers"
	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/sirupsen/logrus"
)

type Service struct {
	settings        ServiceSettings
	serviceCtx      context.Context
	shutdownTimeout time.Duration
	logger          *logrus.Logger
}

func NewService(ctx context.Context, settings ServiceSettings, shutdownTimeout time.Duration, logger *logrus.Logger) *Service {
	return &Service{serviceCtx: ctx, settings: settings, shutdownTimeout: shutdownTimeout, logger: logger}
}

func (service *Service) Start() {
	service.logger.Info("Server service started")

	updHandler := handlers.NewUpdateMetricsHandler(
		handlers.UpdateMetricsUrlPattern,
		storage.NewMemStorage(),
		service.logger,
	)
	updHandler.Handle(http.HandleFunc)

	httpServer := &http.Server{Addr: service.getServerFullAddr(), Handler: nil}

	errChan := make(chan error)
	go func(ch chan error) {
		ch <- httpServer.ListenAndServe()
	}(errChan)

	select {
	case err := <-errChan:
		service.logger.Errorf("Server service exited with err: %v", err)
		return
	case <-service.serviceCtx.Done():
		service.logger.Infof("Server service context canceled")

		ctx, cancel := context.WithTimeout(context.Background(), service.shutdownTimeout)
		defer cancel()

		err := httpServer.Shutdown(ctx)
		if err != nil {
			service.logger.Errorf("Shutdown error: %v", err)
		}
		return
	}
}

func (service *Service) getServerFullAddr() string {
	return fmt.Sprintf("%s:%d", service.settings.serverAddress, service.settings.serverPort)
}
