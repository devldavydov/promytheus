package server

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/devldavydov/promytheus/internal/common/cipher"
	srvgrpc "github.com/devldavydov/promytheus/internal/server/grpc"
	"github.com/devldavydov/promytheus/internal/server/http/handler/metric"
	_middleware "github.com/devldavydov/promytheus/internal/server/http/middleware"
	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/credentials"

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
	// Create storage
	stg, err := service.createStorage(ctx)
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}
	defer stg.Close()

	// Create group for servers
	grp, grpCtx := errgroup.WithContext(ctx)

	// Start HTTP server
	service.startHTTPServer(stg, grp, grpCtx)

	// Start GRPC server
	if service.settings.GRPCAddress != nil {
		service.startGRPCServer(stg, grp, grpCtx)
	}

	return grp.Wait()
}

func (service *Service) loadCryptoPrivKey() (*rsa.PrivateKey, error) {
	if service.settings.CryptoPrivKeyPath == nil {
		return nil, nil
	}
	return cipher.PrivateKeyFromFile(*service.settings.CryptoPrivKeyPath)
}

func (service *Service) createHTTPServer(stg storage.Storage) (*http.Server, error) {
	// Create decryption middleware
	cryptoPrivKey, err := service.loadCryptoPrivKey()
	if err != nil {
		return nil, err
	}
	mdlwrDecr := _middleware.NewDecrpyt(cryptoPrivKey)

	// Create router
	router := chi.NewRouter()
	router.Use(middleware.RealIP, middleware.Logger, middleware.Recoverer, _middleware.Gzip, mdlwrDecr.Handle)

	metric.NewHandler(
		router,
		stg,
		service.settings.HmacKey,
		service.settings.TrustedSubnet,
		service.logger,
	)

	return &http.Server{
			Addr:    service.settings.HTTPAddress.String(),
			Handler: router},
		nil
}

func (service *Service) startHTTPServer(stg storage.Storage, grp *errgroup.Group, grpCtx context.Context) {
	grp.Go(func() error {
		httpServer, err := service.createHTTPServer(stg)
		if err != nil {
			return err
		}

		errChan := make(chan error)
		go func(ch chan error) {
			service.logger.Infof("HTTP service started on [%s]", service.settings.HTTPAddress.String())
			ch <- httpServer.ListenAndServe()
		}(errChan)

		select {
		case err := <-errChan:
			return fmt.Errorf("HTTP service exited with err: %w", err)
		case <-grpCtx.Done():
			service.logger.Infof("HTTP service context canceled")

			ctx, cancel := context.WithTimeout(context.Background(), service.shutdownTimeout)
			defer cancel()

			err := httpServer.Shutdown(ctx)
			if err != nil {
				return fmt.Errorf("HTTP service shutdown err: %w", err)
			}

			service.logger.Info("HTTP service finished")
			return nil
		}
	})
}

func (service *Service) startGRPCServer(stg storage.Storage, grp *errgroup.Group, grpCtx context.Context) {
	grp.Go(func() error {
		var err error

		var tlsCredentials credentials.TransportCredentials
		if service.settings.GRPCServerTLS != nil {
			tlsCredentials, err = service.settings.GRPCServerTLS.Load()
			if err != nil {
				return err
			}
		}

		listen, err := net.Listen("tcp", service.settings.GRPCAddress.String())
		if err != nil {
			return err
		}
		grpcSrv, _ := srvgrpc.NewServer(
			stg,
			service.settings.HmacKey,
			service.settings.TrustedSubnet,
			tlsCredentials,
			service.logger)

		errChan := make(chan error)
		go func(ch chan error) {
			service.logger.Infof("GRPC service started on [%s]", service.settings.GRPCAddress.String())
			ch <- grpcSrv.Serve(listen)
		}(errChan)

		select {
		case err := <-errChan:
			return fmt.Errorf("GRPC service exited with err: %w", err)
		case <-grpCtx.Done():
			service.logger.Infof("GRPC service context canceled")

			grpcSrv.GracefulStop()

			service.logger.Info("GRPC service finished")
			return nil
		}
	})
}

func (service *Service) createStorage(ctx context.Context) (storage.Storage, error) {
	var stg storage.Storage
	var err error

	if service.settings.DatabaseDsn == "" {
		stg, err = storage.NewMemStorage(ctx, service.logger, service.settings.PersistSettings)
	} else {
		stg, err = storage.NewPgStorage(service.settings.DatabaseDsn, service.logger)
	}

	return stg, err
}
