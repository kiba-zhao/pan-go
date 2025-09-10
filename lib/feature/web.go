//go:build !(android || ios)

package feature

import (
	"pan/lib/injection"
	"pan/lib/web"
	"reflect"
)

type WebController interface {
	SetupToWeb(router web.WebRouter) error
}

type WebControllerProvider interface {
	WebControllers() []WebController
}

type WebSubModule interface {
	WebRouteName() string
}

const (
	WebAPIPrefix = "/api/"
)

var _ = (web.WebAppModule)((*SubModule[any, any])(nil))

func (s *SubModule[Module, ParentModule]) SetupToWeb(app web.WebApp) error {
	controllers := getWebControllers(s.specifyModule)
	if len(controllers) <= 0 {
		return nil
	}
	routeName := getWebRouteName(s.specifyModule)
	router := app.Group(WebAPIPrefix + routeName)
	for _, controller := range controllers {
		err := controller.SetupToWeb(router)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SubModule[Module, ParentModule]) Components() []injection.Component {
	controllers := getWebControllers(s.specifyModule)
	if len(controllers) <= 0 {
		return nil
	}

	components := []injection.Component{}
	for _, controller := range controllers {
		t := reflect.TypeOf(controller)
		components = append(components, injection.NewComponentByType(t, controller, injection.ComponentNoneScope))
	}
	return components
}

func getWebControllers(feature interface{}) []WebController {
	if provider, ok := feature.(WebControllerProvider); ok {
		return provider.WebControllers()
	}
	return nil
}

func getWebRouteName(feature interface{}) string {
	var routeName string
	if subModule, ok := feature.(WebSubModule); ok {
		routeName = subModule.WebRouteName()
	}
	return routeName
}
