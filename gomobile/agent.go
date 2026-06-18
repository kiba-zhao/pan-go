//go:build android || ios

package gomobile

import (
	"bytes"
	"errors"
	"io"
	"pan/gobase"
	"pan/pkg/app"
	"pan/pkg/log"
	"pan/pkg/runtime"
	"sync"
)

var errAgentUnavailable = errors.New("gomobile.GoMobileAgent Error: Unavailable")

type GoMobileAgent struct {
	applet app.Applet
	rw     sync.RWMutex

	modules []interface{}
}

func New(cfg SettingsConfig, logger Logger) *GoMobileAgent {
	log.InitDefault(logger)

	module := &stdModule{}

	agent := &GoMobileAgent{}
	agent.modules = []interface{}{
		module,
		gobase.NewSettingsModule(&stdSettingsConfig{cfg: cfg}),
	}
	module.agent = agent

	return agent
}

func (agent *GoMobileAgent) Run() error {

	modules := []interface{}{
		app.New(),
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

func (agent *GoMobileAgent) SetupApplet(applet app.Applet) {
	agent.rw.Lock()
	defer agent.rw.Unlock()
	agent.applet = applet
}

func (agent *GoMobileAgent) Applet() app.Applet {
	agent.rw.RLock()
	defer agent.rw.RUnlock()
	return agent.applet
}

func (agent *GoMobileAgent) Exec(name []byte, in GoMobileReadStream) (GoMobileStream, error) {
	applet := agent.Applet()
	if applet == nil {
		return nil, errAgentUnavailable
	}

	return applet.Exec(name, in)
}

func (agent *GoMobileAgent) Invoke(action []byte, body []byte) ([]byte, error) {
	in := bytes.NewReader(body)
	out, err := agent.Exec(action, in)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(out)
}
