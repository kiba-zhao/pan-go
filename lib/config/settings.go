// Define app settings for the application
package config

import (
	"strconv"
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

	webPort := strconv.FormatUint(uint64(settings.PeerPort), 10)
	settings.WebAddr = "127.0.0.1:" + webPort
	settings.BroadcastAddrs = append(settings.BroadcastAddrs, "224.0.0.2:9001")

	settings.Name = HostName()

	return settings
}
