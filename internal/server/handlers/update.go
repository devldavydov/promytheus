package handlers

import (
	"io"
	"net/http"

	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/sirupsen/logrus"
)

func UpdateMetricsHandler(storage storage.Storage, logger *logrus.Logger) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			io.WriteString(w, "Method Not Allowed\n")
			return
		}

		logger.Debugf("Received update request: %s", req.URL)
	}
}
