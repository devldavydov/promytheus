package handlers

import (
	"io"
	"net/http"

	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/sirupsen/logrus"
)

type UpdateMetricsHandler struct {
	storage storage.Storage
	logger  *logrus.Logger
}

type requestParams struct {
	metricType  string
	metricName  string
	metricValue string
}

func NewUpdateMetricsHandler(storage storage.Storage, logger *logrus.Logger) *UpdateMetricsHandler {
	return &UpdateMetricsHandler{storage: storage, logger: logger}
}

func (handler *UpdateMetricsHandler) Handle() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			io.WriteString(w, "Method Not Allowed\n")
			return
		}

		handler.logger.Debugf("Received update request: %s", req.URL)
	}
}

func (handler *UpdateMetricsHandler) parseRequest(req *http.Request) (requestParams, error) {
	return requestParams{}, nil
}
