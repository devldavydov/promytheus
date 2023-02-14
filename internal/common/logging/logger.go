package logging

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func CreateLogger(logLevel string) (*logrus.Logger, error) {
	logger := logrus.New()
	logLvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, fmt.Errorf("wrong LOG_LEVEL: %w", err)
	}
	logger.SetLevel(logLvl)
	return logger, nil
}
