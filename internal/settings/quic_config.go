package settings

import (
	"crypto/tls"
	"pan/internal/quic"
	"time"
)

type stdQuicConfig struct {
	security      SecurityConfig
	port          uint16
	addrs         []string
	certificate   tls.Certificate
	dialThreshold uint16
	dialTimeout   time.Duration
}

func newQuicConfig(settings *Settings, security SecurityConfig, netIfaces []NetInterface) quic.QuicConfig {
	cfg := &stdQuicConfig{}
	cfg.port = settings.PeerPort
	cfg.security = security

	if settings.Enabled {
		initQuicConfigWithNetInterfaces(cfg, netIfaces)
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
	return 1000
}
func (cfg *stdQuicConfig) DialTimeout() time.Duration {
	return time.Second * 60
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
