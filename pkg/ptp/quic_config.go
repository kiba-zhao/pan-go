package ptp

import (
	"crypto"
	"crypto/tls"
	"pan/pkg/config"
)

type QuicConfig interface {
	Port() uint16
	Addrs() []string
	Certificate() tls.Certificate
	PeerID() PeerID
	PrivateKey() crypto.PrivateKey
	BroadcastEnabled() bool
}

type QuicConfigListener config.ConfigurerListener[QuicConfig]
type QuicConfigurer config.Configurer[QuicConfig]
