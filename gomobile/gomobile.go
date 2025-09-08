//go:build android || ios

package gomobile

import (
	"pan/gobase"
	"pan/lib/log"
)

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
