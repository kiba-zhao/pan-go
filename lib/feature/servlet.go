//go:build android || ios

package feature

import (
	"pan/lib/injection"
	"pan/lib/servlet"
	"reflect"
)

type ServletHandler interface {
	SetupToServlet(router servlet.ServletRouter) error
}

type ServletHandlerProvider interface {
	ServletHandlers() []ServletHandler
}

type ServletSubModule interface {
	ServletRouteName() []byte
}

var _ = (servlet.ServletModule)((*SubModule[any, any])(nil))

func (s *SubModule[Module, ParentModule]) SetupToServlet(app servlet.ServletApp) error {
	handlers := getServletHandlers(s.specifyModule)
	if len(handlers) <= 0 {
		return nil
	}
	routeName := getServletRouteName(s.specifyModule)
	router := app.Route(routeName)
	for _, handler := range handlers {
		err := handler.SetupToServlet(router)
		if err != nil {
			return err
		}
	}
	return nil
}

var _ = (injection.ComponentProvider)((*SubModule[any, any])(nil))

func (s *SubModule[Module, ParentModule]) Components() []injection.Component {

	handlers := getServletHandlers(s.specifyModule)
	if len(handlers) <= 0 {
		return nil
	}
	components := []injection.Component{}
	for _, handler := range handlers {
		t := reflect.TypeOf(handler)
		components = append(components, injection.NewComponentByType(t, handler, injection.ComponentNoneScope))
	}
	return components
}

func getServletHandlers(feature interface{}) []ServletHandler {
	var handlers []ServletHandler
	if provider, ok := feature.(ServletHandlerProvider); ok {
		handlers = provider.ServletHandlers()
	}
	return handlers
}

func getServletRouteName(feature interface{}) []byte {
	var routeName []byte
	if subModule, ok := feature.(ServletSubModule); ok {
		routeName = subModule.ServletRouteName()
	}
	return routeName
}
