package bootstrap

import "pan/runtime"

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
	err := e.DeferEngine.bootstrap()
	if err == nil {
		err = e.ReadyEngine.bootstrap()
	}
	return err
}

func (e *engine) Components() []Component {
	return []Component{
		NewComponent(e, ComponentNoneScope),
	}
}
