package settings

import (
	"crypto/tls"
	"pan/lib/peer"
	"pan/lib/quic"
	"time"
)

type stdQuicConfig struct {
	security      peer.PeerSecurity
	port          uint16
	addrs         []string
	certificate   tls.Certificate
	dialThreshold uint16
	dialTimeout   time.Duration
}

func parseQuicConfig(settings Settings, security peer.PeerSecurity) *stdQuicConfig {
	cfg := &stdQuicConfig{}
	cfg.port = settings.PeerPort
	cfg.dialThreshold = 1000
	cfg.dialTimeout = time.Second * 60
	cfg.security = security

	if settings.Enabled {
		cfg.addrs = append(cfg.addrs, "")
	}
	return cfg
}

var _ = (quic.QuicConfig)((*stdQuicConfig)(nil))

func (cfg *stdQuicConfig) Port() uint16 {
	return cfg.port
}
func (cfg *stdQuicConfig) Addrs() []string {
	return cfg.addrs
}
func (cfg *stdQuicConfig) Certificate() tls.Certificate {
	return cfg.security.Certificate()
}
func (cfg *stdQuicConfig) DialThreshold() uint16 {
	return cfg.dialThreshold
}
func (cfg *stdQuicConfig) DialTimeout() time.Duration {
	return cfg.dialTimeout
}
