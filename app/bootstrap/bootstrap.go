// Package bootstrap provides the bootstrap engine
//
// The bootstrap engine is used to manage the bootstrap of the application
package bootstrap

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"pan/app/injection"
	"pan/runtime"
	"syscall"
)

var ErrBootstrapExit = errors.New("bootstrap.engine Error: bootstrap exit")

// New returns the engine for bootstrap
func New() interface{} {
	return runtime.NewModule(injection.New(), &readyEngine{}, &deferEngine{})
}

// Bootstrap returns the bootstrap engine
func Bootstrap() interface{} {
	return &engine{}
}

type engine struct {
	ReadyEngine *readyEngine
	DeferEngine *deferEngine
}

// Init initializes the bootstrap engine with the provided registry.
// It sets up signal handling for graceful shutdown on SIGINT and SIGTERM.
// The function creates a context with cancellation and triggers the bootstrap
// process for both the deferEngine and readyEngine. If any errors occur during
// the bootstrap process, the error is returned.

func (e *engine) Init(registry runtime.Registry) error {

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancelCause(context.Background())
	go func() {
		<-signals
		cancel(ErrBootstrapExit)
	}()

	err := e.DeferEngine.bootstrap(ctx)
	if err == nil {
		err = e.ReadyEngine.bootstrap(ctx)
	}

	return err
}

// Components returns a slice of injection.Component that are provided by the
// bootstrap engine. This is just the engine itself, and the scope is set to
// ComponentNoneScope as the bootstrap engine is not injected into any other
// components.
func (e *engine) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(e, injection.ComponentNoneScope),
	}
}
