// Define app settings for the application
package config

import (
	"os"
	"pan/lib/net"
)

type Settings struct {
	Name             string   `json:"name" form:"name"`
	WebAddress       []string `json:"webAddress" form:"webAddress"`
	PeerAddress      []string `json:"peerAddress" form:"peerAddress"`
	BroadcastAddress []string `json:"broadcastAddress" form:"broadcastAddress"`
	PublicAddress    []string `json:"publicAddress" form:"publicAddress"`
	GuardEnabled     bool     `json:"guardEnabled" form:"guardEnabled"`
	GuardAccess      bool     `json:"guardAccess" form:"guardAccess"`
	DiscoveryServer  []string `json:"discoveryServer" form:"discoveryServer"`
}

func newDefaultSettings(cfg AppConfig) AppSettings {

	settings := &Settings{}
	addrStat, err := net.StatAddr()
	if err == nil {
		if addrStat.IPv6Enabled {
			settings.WebAddress = append(settings.WebAddress, "[::1]:9002")
			settings.PeerAddress = append(settings.PeerAddress, "[::]:9000")
			settings.BroadcastAddress = append(settings.BroadcastAddress, "[FF02::1]:9100")
		}
		if addrStat.IPv4Enabled {
			settings.WebAddress = append(settings.WebAddress, "127.0.0.1:9002")
			if !addrStat.IPv6Enabled {
				settings.PeerAddress = append(settings.PeerAddress, "0.0.0.0:9000")
			}
			settings.BroadcastAddress = append(settings.BroadcastAddress, "224.0.0.2:9100")
		}
		// Temporary annotation, awaiting completion of broadcast optimization
		// if addrStat.IPv6GlobalEnabled {
		// 	settings.BroadcastAddress = append(settings.BroadcastAddress, "[FF0E::1]:9100")
		// }
		if addrStat.IPv4GlobalEnabled {
			settings.BroadcastAddress = append(settings.BroadcastAddress, "224.0.1.1:9100")
		}
	}

	settings.Name = generateName()
	settings.PublicAddress = settings.PeerAddress
	settings.GuardEnabled = true
	settings.GuardAccess = true

	return settings
}

func generateName() string {
	name, err := os.Hostname()
	if err == nil {
		return name
	}

	return "pan-go"
}
