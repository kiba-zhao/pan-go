package feature

import (
	"errors"
	"pan/lib/serlvet"
)

type SerlvetHandler interface {
	SetupToSerlvet(router serlvet.SerlvetRouter) error
}

type SerlvetHandlerProvider interface {
	SerlvetHandlers() []SerlvetHandler
}

var _ = (serlvet.SerlvetModule)((*stdFeatureModule)(nil))

func (module *stdFeatureModule) SetupToSerlvet(app serlvet.SerlvetApp) error {

	var errs []error
	featureHelper := module.featureHelper
	handlers := getSerlvetHandlers(featureHelper.feature)
	if len(handlers) <= 0 {
		return nil
	}

	router := app.Route(featureHelper.SerlvetScope())
	for _, handler := range handlers {
		err := handler.SetupToSerlvet(router)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func getSerlvetHandlers(feature interface{}) []SerlvetHandler {
	var handlers []SerlvetHandler
	if provider, ok := feature.(SerlvetHandlerProvider); ok {
		handlers = provider.SerlvetHandlers()
	}
	return handlers
}
