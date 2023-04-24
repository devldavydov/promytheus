// Package log provides functions for logging.
package log

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// NewLogger creates new logger with specified arguments.
func NewLogger(logLevel string, logFile string) (*logrus.Logger, io.Closer, error) {
	logger := logrus.New()
	logLvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, nil, fmt.Errorf("wrong LOG_LEVEL: %w", err)
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}

	logger.SetLevel(logLvl)
	logger.SetOutput(file)
	return logger, file, nil
}
