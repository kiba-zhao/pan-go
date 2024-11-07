package bootstrap

import (
	"context"
	"os"
	"os/signal"
	"pan/app/constant"
	"pan/runtime"
	"syscall"
)

func New() interface{} {
	return runtime.NewModule(&injectEngine{}, &readyEngine{}, &deferEngine{})
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
		cancel(constant.ErrExit)
	}()

	err := e.DeferEngine.bootstrap(ctx)
	if err == nil {
		err = e.ReadyEngine.bootstrap(ctx)
	}

	return err
}

func (e *engine) Components() []Component {
	return []Component{
		NewComponent(e, ComponentNoneScope),
	}
}
