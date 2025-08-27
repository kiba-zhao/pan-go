package broadcast

import (
	"crypto"
	"pan/lib/config"
	"pan/lib/peer"
)

type BroadcastConfig interface {
	Interfaces() []BroadcastInterface

	Addrs() []string
	MTU() int
	IPv6Enabled() bool

	DeliverMaxSize() int
	IPv6ZoneList() []string

	PeerID() peer.PeerID
	PrivateKey() crypto.PrivateKey
}

type BroadcastConfigListener = config.ConfigurerListener[BroadcastConfig]
type BroadcastConfigurer = config.Configurer[BroadcastConfig]
