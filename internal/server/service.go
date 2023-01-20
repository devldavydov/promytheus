package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/devldavydov/promytheus/internal/server/handlers"
	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

type Service struct {
	settings        ServiceSettings
	shutdownTimeout time.Duration
	logger          *logrus.Logger
}

func NewService(settings ServiceSettings, shutdownTimeout time.Duration, logger *logrus.Logger) *Service {
	return &Service{settings: settings, shutdownTimeout: shutdownTimeout, logger: logger}
}

func (service *Service) Start(ctx context.Context) {
	service.logger.Info("Server service started")

	metricsHandler := handlers.NewMetricsHandler(
		storage.NewMemStorage(),
		service.logger,
	)
	r := handlers.NewRouter(metricsHandler, middleware.RealIP, middleware.Logger, middleware.Recoverer)

	httpServer := &http.Server{Addr: service.getServerFullAddr(), Handler: r}

	errChan := make(chan error)
	go func(ch chan error) {
		ch <- httpServer.ListenAndServe()
	}(errChan)

	select {
	case err := <-errChan:
		service.logger.Errorf("Server service exited with err: %v", err)
		return
	case <-ctx.Done():
		service.logger.Infof("Server service context canceled")

		ctx, cancel := context.WithTimeout(context.Background(), service.shutdownTimeout)
		defer cancel()

		err := httpServer.Shutdown(ctx)
		if err != nil {
			service.logger.Errorf("Shutdown error: %v", err)
		}

		service.logger.Info("Server service finished")
		return
	}
}

func (service *Service) getServerFullAddr() string {
	return fmt.Sprintf("%s:%d", service.settings.serverAddress, service.settings.serverPort)
}
