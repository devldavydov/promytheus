package server

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"

	"github.com/devldavydov/promytheus/internal/common/cipher"
	"github.com/devldavydov/promytheus/internal/server/handler/metric"
	_middleware "github.com/devldavydov/promytheus/internal/server/middleware"
	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

type Service struct {
	logger          *logrus.Logger
	settings        ServiceSettings
	shutdownTimeout time.Duration
}

func NewService(settings ServiceSettings, shutdownTimeout time.Duration, logger *logrus.Logger) *Service {
	return &Service{settings: settings, shutdownTimeout: shutdownTimeout, logger: logger}
}

func (service *Service) Start(ctx context.Context) error {
	service.logger.Infof("Server service started on [%s:%d]", service.settings.ServerAddress, service.settings.ServerPort)

	// Create decryption middleware
	cryptoPrivKey, err := service.loadCryptoPrivKey()
	if err != nil {
		return err
	}
	mdlwrDecr := _middleware.NewDecrpyt(cryptoPrivKey)

	// Create router
	router := chi.NewRouter()
	router.Use(middleware.RealIP, middleware.Logger, middleware.Recoverer, _middleware.Gzip, mdlwrDecr.Handle)

	var stg storage.Storage

	if service.settings.DatabaseDsn == "" {
		stg, err = storage.NewMemStorage(ctx, service.logger, service.settings.PersistSettings)
	} else {
		stg, err = storage.NewPgStorage(service.settings.DatabaseDsn, service.logger)
	}

	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}
	defer stg.Close()

	metric.NewHandler(
		router,
		stg,
		service.settings.HmacKey,
		service.logger,
	)

	httpServer := &http.Server{Addr: service.getServerFullAddr(), Handler: router}

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

func (service *Service) loadCryptoPrivKey() (*rsa.PrivateKey, error) {
	if service.settings.CryptoPrivKeyPath == nil {
		return nil, nil
	}
	return cipher.PrivateKeyFromFile(*service.settings.CryptoPrivKeyPath)
}
