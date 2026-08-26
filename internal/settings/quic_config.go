package settings

import (
	"crypto"
	"crypto/tls"
	"pan/pkg/ptp"
	"time"
)

type stdQuicConfig struct {
	security         SecurityConfig
	port             uint16
	addrs            []string
	certificate      tls.Certificate
	dialThreshold    uint16
	dialTimeout      time.Duration
	broadcastEnabled bool
}

func newQuicConfig(deviceNetwork *DeviceNetwork, security SecurityConfig, addrs []string) ptp.QuicConfig {
	cfg := &stdQuicConfig{}
	cfg.port = deviceNetwork.Port
	cfg.security = security
	cfg.broadcastEnabled = deviceNetwork.BroadcastEnabled

	if deviceNetwork.Enabled {
		cfg.addrs = addrs
	}

	return cfg
}

var _ = (ptp.QuicConfig)((*stdQuicConfig)(nil))

func (cfg *stdQuicConfig) Port() uint16 {
	return cfg.port
}
func (cfg *stdQuicConfig) Addrs() []string {
	return cfg.addrs
}
func (cfg *stdQuicConfig) Certificate() tls.Certificate {
	return cfg.security.Certificate()
}

func (cfg *stdQuicConfig) PeerID() ptp.PeerID {
	return cfg.security.PeerID()
}

func (cfg *stdQuicConfig) PrivateKey() crypto.PrivateKey {
	return cfg.security.PrivateKey()
}

func (cfg *stdQuicConfig) BroadcastEnabled() bool {
	return cfg.broadcastEnabled
}
