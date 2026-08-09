package gobase

import (
	"pan/internal/settings"
	"pan/pkg/bootstrap"
	"pan/pkg/ptp"
	"pan/pkg/repository"
	"pan/pkg/runtime"
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
		ptp.New(),
	}

	if len(module.modules) > 0 {
		modules = append(modules, module.modules...)
	}

	modules = append(
		modules,
		settings.New(),
		// appinfo.New(),
		// user.New(),
		bootstrap.Bootstrap(),
	)
	return modules
}
