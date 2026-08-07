//go:build !(android || ios)

package settings

import (
	"net"
	"pan/pkg/web"
	"strconv"
)

type stdWebConfig struct {
	settings       *HostSettings
	settingsConfig SettingsConfig
}

func newWebConfig(settings *HostSettings, settingsConfig SettingsConfig) web.WebConfig {
	cfg := &stdWebConfig{}
	cfg.settings = settings
	cfg.settingsConfig = settingsConfig
	return cfg
}

var IPv4LocalHost = net.IP{127, 0, 0, 1}
var IPv6LocalHost = net.IPv6loopback
var _ = (web.WebConfig)((*stdWebConfig)(nil))

func (cfg *stdWebConfig) Addr() string {
	settings := cfg.settings
	if !settings.WebEnabled {
		return ""
	}

	ipv6Enabled := cfg.settingsConfig.IPv6Enabled()
	webPort := settings.WebPort
	var ip net.IP
	if settings.LocalHostOnly {
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
