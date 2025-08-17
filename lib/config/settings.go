// Define app settings for the application
package config

import (
	"pan/lib/net"
)

type AppSettings = *Settings
type AppConfig = Config[AppSettings]
type AppConfigListener = ConfigListener[AppSettings]

type Settings struct {
	Name           string   `json:"name" form:"name"`
	PeerPort       uint16   `json:"peerPort" form:"peerPort"`
	BroadcastAddrs []string `json:"broadcastAddrs" form:"broadcastAddrs"`
	PublicAddrs    []string `json:"publicAddrs" form:"publicAddrs"`
	Enabled        bool     `json:"enabled" form:"enabled"`
	WebAddr        string   `json:"webAddr" form:"webAddr"`
}

func NewDefaultSettings() AppSettings {

	settings := &Settings{}
	settings.PeerPort = 9000
	settings.Enabled = true

	addrStat, err := net.StatAddr()
	if err == nil {
		if addrStat.IPv6Enabled {
			settings.WebAddr = "[::1]:9002"
			settings.BroadcastAddrs = append(settings.BroadcastAddrs, "[FF02::1]:9100")
		}
		if addrStat.IPv4Enabled {
			settings.WebAddr = "127.0.0.1:9002"
			settings.BroadcastAddrs = append(settings.BroadcastAddrs, "224.0.0.2:9100")
		}
		// Temporary annotation, awaiting completion of broadcast optimization
		// if addrStat.IPv6GlobalEnabled {
		// 	settings.BroadcastAddress = append(settings.BroadcastAddress, "[FF0E::1]:9100")
		// }
		if addrStat.IPv4GlobalEnabled {
			settings.BroadcastAddrs = append(settings.BroadcastAddrs, "224.0.1.1:9100")
		}
	}

	settings.Name = HostName()

	return settings
}
