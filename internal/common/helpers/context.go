package helpers

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func CreateContextWithSignalHadler(ctxParent context.Context) context.Context {
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
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx
}
