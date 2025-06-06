package feature

import (
	"errors"
	"pan/lib/web"
)

type WebController interface {
	SetupToWeb(router web.WebRouter) error
}

type WebControllerProvider interface {
	WebControllers() []WebController
}

var _ = (web.WebAppModule)((*stdFeatureModule)(nil))

func (module *stdFeatureModule) SetupToWeb(app web.WebApp) error {

	var errs []error
	featureHelper := module.featureHelper
	controllers := getWebControllers(featureHelper.feature)
	if len(controllers) <= 0 {
		return nil
	}

	router := app.Group(featureHelper.WebScope())
	for _, controller := range controllers {
		err := controller.SetupToWeb(router)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func getWebControllers(feature interface{}) []WebController {
	if provider, ok := feature.(WebControllerProvider); ok {
		return provider.WebControllers()
	}
	return nil
}
