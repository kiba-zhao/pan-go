package gomobile

import (
	"pan/features/app"
	"pan/features/extfs"
	"pan/lib/runtime"
	libSerlvet "pan/lib/serlvet"
	"sync"
)

type GoMobileAgent struct {
	serlvet *libSerlvet.Serlvet
	rw      sync.RWMutex

	module *stdModule
}

func (agent *GoMobileAgent) Run() error {

	serlvetModule := libSerlvet.New()
	engine, err := runtime.New(
		app.New(serlvetModule, agent.module),
		extfs.New(),
		app.Bootstrap(),
	)

	if err == nil {
		ctx := runtime.NewContext()
		err = engine.Bootstrap(ctx)
	}
	return err

}

func (agent *GoMobileAgent) Terminate() error {
	return runtime.AbortContext()
}

func getAgentSerlvet(agent *GoMobileAgent) *libSerlvet.Serlvet {
	agent.rw.RLock()
	defer agent.rw.RUnlock()
	return agent.serlvet
}

func setAgentSerlvet(agent *GoMobileAgent, serlvet *libSerlvet.Serlvet) {
	agent.rw.Lock()
	defer agent.rw.Unlock()
	agent.serlvet = serlvet
}
