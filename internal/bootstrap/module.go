// Package bootstrap provides the bootstrap engine
//
// The bootstrap engine is used to manage the bootstrap of the application
package bootstrap

import (
	"context"
	"errors"
	"pan/internal/injection"
	"pan/internal/runtime"
)

var ErrBootstrapExit = errors.New("bootstrap.engine Error: bootstrap exit")

// New returns the engine for bootstrap
func New() interface{} {
	return runtime.NewModule(injection.New(), &readyEngine{}, &deferEngine{}, &destroyEngine{})
}

// Bootstrap returns the bootstrap engine
func Bootstrap() interface{} {
	return &stdBootstrapModule{}
}

type stdBootstrapModule struct {
	ReadyEngine   *readyEngine
	DeferEngine   *deferEngine
	DestroyEngine *destroyEngine
}

var _ = (runtime.InitializeModule)((*stdBootstrapModule)(nil))

func (e *stdBootstrapModule) Init(ctx context.Context, registry runtime.Registry) error {

	err := e.DeferEngine.bootstrap(ctx)
	if err == nil {
		err = e.ReadyEngine.bootstrap(ctx)
	}

	if err == nil {
		err = e.DestroyEngine.bootstrap(ctx)
	}
	return err
}

var _ = (injection.ComponentProvider)((*stdBootstrapModule)(nil))

func (e *stdBootstrapModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(e, injection.ComponentNoneScope),
	}
}
