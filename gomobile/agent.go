//go:build android || ios

package gomobile

import (
	"bytes"
	"context"
	"errors"
	"io"
	"pan/gobase"
	"pan/lib/runtime"
	libServlet "pan/lib/servlet"
	"sync"
)

var errAgentUnavailable = errors.New("GoMobileAgent Error: Unavailable")

type GoMobileAgent struct {
	servlet libServlet.Servlet
	rw      sync.RWMutex

	modules []interface{}
}

func (agent *GoMobileAgent) Run() error {

	modules := []interface{}{
		libServlet.New(),
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
	servlet := getAgentServlet(agent)
	if servlet == nil {
		return nil, errAgentUnavailable
	}

	ctx := context.Background()
	return servlet.Do(ctx, action, in)
}

func (agent *GoMobileAgent) Exec(action []byte, body []byte) ([]byte, error) {
	in := bytes.NewReader(body)
	out, err := agent.Do(action, in)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(out)
}

func getAgentServlet(agent *GoMobileAgent) libServlet.Servlet {
	agent.rw.RLock()
	defer agent.rw.RUnlock()
	return agent.servlet
}

func setupAgentServlet(agent *GoMobileAgent, servlet libServlet.Servlet) {
	agent.rw.Lock()
	defer agent.rw.Unlock()
	agent.servlet = servlet
}
