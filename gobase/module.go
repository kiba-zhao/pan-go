package gobase

import (
	"pan/internal/appinfo"
	"pan/internal/bootstrap"
	"pan/internal/broadcast"
	"pan/internal/peer"
	"pan/internal/quic"
	"pan/internal/repository"
	"pan/internal/runtime"
	"pan/internal/settings"
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
		appinfo.New(),
		bootstrap.Bootstrap(),
	)
	return modules
}
