package test

import (
	"pan/app/services"
	"pan/core"
)

type AppEnabledModule interface {
	core.AppModule
	services.EnabledModule
}
