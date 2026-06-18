package net

import (
	"pan/pkg/config"
)

type BroadcastInterface struct {
	Addr string
	MTU  int
}

type BroadcastConfig interface {
	Interfaces() []BroadcastInterface

	Addrs() []string
	MTU() int
	IPv6ZoneList() []string

	DeliverMaxSize() int
}

type BroadcastConfigListener = config.ConfigurerListener[BroadcastConfig]
type BroadcastConfigurer = config.Configurer[BroadcastConfig]
