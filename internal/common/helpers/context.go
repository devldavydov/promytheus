package helpers

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func CreateContextWithSignalHadler(ctxParent context.Context, logger *logrus.Logger) context.Context {
	ctx, cancel := context.WithCancel(ctxParent)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer func() {
			signal.Stop(c)
			cancel()
		}()

		select {
		case <-c:
			logger.Info("Received signal, stoppping...")
			cancel()
		case <-ctx.Done():
			logger.Info("Context canceled, stopping...")
		}
	}()

	return ctx
}
