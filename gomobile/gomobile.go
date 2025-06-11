package gomobile

import "pan/lib/log"

func New(settings Settings, logger Logger) *GoMobileAgent {
	log.InitDefault(logger)
	initWithSettings(settings, logger)

	agent := &GoMobileAgent{}

	module := &stdModule{}
	agent.module = module
	module.agent = agent

	return agent
}
