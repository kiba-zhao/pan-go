package web

import "pan/lib/config"

type WebConfig interface {
	Addr() string
}

type WebConfigListener = config.ConfigurerListener[WebConfig]
type WebConfigurer = config.Configurer[WebConfig]
