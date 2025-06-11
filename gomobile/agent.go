package gomobile

import (
	"bytes"
	"context"
	"errors"
	"io"
	"pan/features/app"
	"pan/features/extfs"
	"pan/lib/runtime"
	libSerlvet "pan/lib/serlvet"
	"sync"
)

var errAgentUnavailable = errors.New("GoMobileAgent Error: Unavailable")

type GoMobileAgent struct {
	serlvet libSerlvet.Serlvet
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

func (agent *GoMobileAgent) Do(action []byte, in GoMobileReadStream) (GoMobileStream, error) {
	serlvet := getAgentSerlvet(agent)
	if serlvet == nil {
		return nil, errAgentUnavailable
	}

	ctx := context.Background()
	return serlvet.Do(ctx, action, in)
}

func (agent *GoMobileAgent) Exec(action []byte, body []byte) ([]byte, error) {
	in := bytes.NewReader(body)
	out, err := agent.Do(action, in)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(out)
}

func getAgentSerlvet(agent *GoMobileAgent) libSerlvet.Serlvet {
	agent.rw.RLock()
	defer agent.rw.RUnlock()
	return agent.serlvet
}

func setAgentSerlvet(agent *GoMobileAgent, serlvet libSerlvet.Serlvet) {
	agent.rw.Lock()
	defer agent.rw.Unlock()
	agent.serlvet = serlvet
}
