package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/devldavydov/promytheus/internal/server/handler"
	_middleware "github.com/devldavydov/promytheus/internal/server/middleware"
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

func (service *Service) Start(ctx context.Context) error {
	service.logger.Infof("Server service started on [%s:%d]", service.settings.ServerAddress, service.settings.ServerPort)

	memStorage, err := storage.NewMemStorage(ctx, service.logger, service.settings.PersistSettings)
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}

	metricsHandler := handler.NewMetricsHandler(
		memStorage,
		service.settings.HmacKey,
		service.logger,
	)
	r := handler.NewRouter(metricsHandler, middleware.RealIP, middleware.Logger, middleware.Recoverer, _middleware.Gzip)

	httpServer := &http.Server{Addr: service.getServerFullAddr(), Handler: r}

	errChan := make(chan error)
	go func(ch chan error) {
		ch <- httpServer.ListenAndServe()
	}(errChan)

	select {
	case err := <-errChan:
		return fmt.Errorf("server service exited with err: %w", err)
	case <-ctx.Done():
		service.logger.Infof("Server service context canceled")

		ctx, cancel := context.WithTimeout(context.Background(), service.shutdownTimeout)
		defer cancel()

		err := httpServer.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("server service shutdown err: %w", err)
		}

		service.logger.Info("Server service finished")
		return nil
	}
}

func (service *Service) getServerFullAddr() string {
	return fmt.Sprintf("%s:%d", service.settings.ServerAddress, service.settings.ServerPort)
}
