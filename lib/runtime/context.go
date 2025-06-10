package runtime

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
)

var errSignalAbort = errors.New("runtime.context Error: Signal Abort")
var errContextUnavailable = errors.New("runtime.context Error: Context Unavailable")

func EnsureContext(ctx context.Context) error {
	if ctx == nil {
		return errContextUnavailable
	}

	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	default:
	}
	return err
}

func IsContextUnavailable(err error) bool {
	return errors.Is(err, errContextUnavailable)
}

func NewContext() context.Context {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancelCause(context.Background())
	go func() {
		<-signals
		cancel(errSignalAbort)
	}()
	return ctx
}

func IsAbort(err error) bool {
	return errors.Is(err, errSignalAbort)
}

func AbortContext() error {
	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}
	return process.Signal(syscall.SIGINT)
}
