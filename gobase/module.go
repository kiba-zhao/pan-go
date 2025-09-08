package gobase

import (
	"pan/features/settings"
	"pan/lib/bootstrap"
	"pan/lib/broadcast"
	"pan/lib/peer"
	"pan/lib/quic"
	"pan/lib/repository"
	"pan/lib/runtime"
)

func New(modules ...interface{}) interface{} {
	module := &stdModule{}
	module.modules = modules

	return module
}

type stdModule struct {
	modules []interface{}
}

var _ = (runtime.ProviderModule)((*stdModule)(nil))

func (module *stdModule) Modules() []interface{} {

	modules := []interface{}{
		bootstrap.New(),
		repository.New(),
		peer.New(),
		broadcast.New(),
		quic.New(),
	}

	if len(module.modules) > 0 {
		modules = append(modules, module.modules...)
	}

	modules = append(
		modules,
		settings.New(),
		bootstrap.Bootstrap(),
	)
	return modules
}
