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

func New() interface{} {
	return runtime.NewModule(injection.New(), &readyEngine{}, &deferEngine{})
}

func Bootstrap() interface{} {
	return &engine{}
}

type engine struct {
	ReadyEngine *readyEngine
	DeferEngine *deferEngine
}

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

func (e *engine) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(e, injection.ComponentNoneScope),
	}
}
