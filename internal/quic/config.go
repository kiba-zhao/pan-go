package quic

import (
	"crypto/tls"
	"pan/internal/config"
	"time"
)

type QuicConfig interface {
	Port() uint16
	Addrs() []string
	Certificate() tls.Certificate
	DialThreshold() uint16
	DialTimeout() time.Duration
}

type QuicConfigListener config.ConfigurerListener[QuicConfig]
type QuicConfigurer config.Configurer[QuicConfig]
