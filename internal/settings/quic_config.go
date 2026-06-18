package settings

import (
	"crypto"
	"crypto/tls"
	"pan/pkg/net"
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

func newQuicConfig(settings *Settings, security SecurityConfig, netIfaces []NetInterface) net.QuicConfig {
	cfg := &stdQuicConfig{}
	cfg.port = settings.PeerPort
	cfg.security = security
	cfg.broadcastEnabled = settings.BroadcastEnabled

	if settings.Enabled {
		initQuicConfigWithNetInterfaces(cfg, netIfaces)
	}
	return cfg
}

var _ = (net.QuicConfig)((*stdQuicConfig)(nil))

func (cfg *stdQuicConfig) Port() uint16 {
	return cfg.port
}
func (cfg *stdQuicConfig) Addrs() []string {
	return cfg.addrs
}
func (cfg *stdQuicConfig) Certificate() tls.Certificate {
	return cfg.security.Certificate()
}

func (cfg *stdQuicConfig) PeerID() net.PeerID {
	return cfg.security.PeerID()
}

func (cfg *stdQuicConfig) PrivateKey() crypto.PrivateKey {
	return cfg.security.PrivateKey()
}

func (cfg *stdQuicConfig) BroadcastEnabled() bool {
	return cfg.broadcastEnabled
}

func initQuicConfigWithNetInterfaces(cfg *stdQuicConfig, netIfaces []NetInterface) {
	if len(netIfaces) <= 0 {
		cfg.addrs = append(cfg.addrs, "")
	} else {
		for _, iface := range netIfaces {
			cfg.addrs = append(cfg.addrs, iface.Addr)
		}
	}
}
