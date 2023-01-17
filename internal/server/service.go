package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

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

	metricsHandler := func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			io.WriteString(w, "Method Not Allowed\n")
			return
		}

		fmt.Println(req.URL)
	}

	http.HandleFunc("/", metricsHandler)
	httpServer := &http.Server{Addr: fmt.Sprintf("%s:%d", service.settings.serverAddress, service.settings.serverPort), Handler: nil}

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
