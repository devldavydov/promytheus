package helpers

import (
	"context"
	"os/signal"
	"syscall"
)

func CreateContextWithSignalHadler(parent context.Context) context.Context {
	ctx, stop := signal.NotifyContext(parent, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer stop()
		<-ctx.Done()
	}()

	return ctx
}
