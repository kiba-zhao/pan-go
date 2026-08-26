//go:build !(android || ios)

package settings

import (
	"net"
	"pan/pkg/web"
	"strconv"
)

type stdWebConfig struct {
	webHost        *WebHost
	settingsConfig SettingsConfig
}

func newWebConfig(webHost *WebHost, settingsConfig SettingsConfig) web.WebConfig {
	cfg := &stdWebConfig{}
	cfg.webHost = webHost
	cfg.settingsConfig = settingsConfig
	return cfg
}

var IPv4LocalHost = net.IP{127, 0, 0, 1}
var IPv6LocalHost = net.IPv6loopback
var _ = (web.WebConfig)((*stdWebConfig)(nil))

func (cfg *stdWebConfig) Addr() string {
	webHost := cfg.webHost
	if !webHost.WebEnabled {
		return ""
	}

	ipv6Enabled := cfg.settingsConfig.IPv6Enabled()
	webPort := webHost.WebPort
	var ip net.IP
	if webHost.LocalHostOnly {
		if ipv6Enabled {
			ip = IPv6LocalHost
		} else {
			ip = IPv4LocalHost
		}
	} else {
		if ipv6Enabled {
			ip = net.IPv6zero
		} else {
			ip = net.IPv4zero
		}
	}

	return net.JoinHostPort(ip.String(), strconv.Itoa(int(webPort)))
}
