package broadcast

import (
	"crypto"
	"pan/internal/config"
	"pan/internal/peer"
)

type BroadcastConfig interface {
	Interfaces() []BroadcastInterface

	Addrs() []string
	MTU() int
	IPv6ZoneList() []string

	DeliverMaxSize() int

	PeerID() peer.PeerID
	PrivateKey() crypto.PrivateKey
}

type BroadcastConfigListener = config.ConfigurerListener[BroadcastConfig]
type BroadcastConfigurer = config.Configurer[BroadcastConfig]
