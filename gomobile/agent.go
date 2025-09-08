//go:build android || ios

package gomobile

import (
	"bytes"
	"context"
	"errors"
	"io"
	"pan/gobase"
	"pan/lib/runtime"
	libSerlvet "pan/lib/serlvet"
	"sync"
)

var errAgentUnavailable = errors.New("GoMobileAgent Error: Unavailable")

type GoMobileAgent struct {
	serlvet libSerlvet.Serlvet
	rw      sync.RWMutex

	modules []interface{}
}

func (agent *GoMobileAgent) Run() error {

	modules := []interface{}{
		libSerlvet.New(),
	}
	modules = append(modules, agent.modules...)
	module := gobase.New(modules...)

	engine, err := runtime.New(module)
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

func setupAgentSerlvet(agent *GoMobileAgent, serlvet libSerlvet.Serlvet) {
	agent.rw.Lock()
	defer agent.rw.Unlock()
	agent.serlvet = serlvet
}
