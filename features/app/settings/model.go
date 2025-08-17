package settings

import "pan/lib/config"

type AppSettings struct {
	config.Settings
	PeerID     string `json:"peerId" form:"peerId"`
	ConfigPath string `json:"cfgPath" form:"cfgPath"`
}

type AppSettingsFields struct {
	Name           string   `form:"name" json:"name"  binding:"omitempty"`
	PeerPort       *uint16  `form:"peerPort" json:"peerPort"  binding:"omitempty"`
	BroadcastAddrs []string `form:"broadcastAddrs" json:"broadcastAddrs"  binding:"omitempty"`
	PublicAddrs    []string `form:"publicAddrs" json:"publicAddrs"  binding:"omitempty"`
	Enabled        *bool    `form:"enabled" json:"enabled"  binding:"omitempty"`
	WebAddr        *string  `form:"webAddr" json:"webAddr"  binding:"omitempty"`
}
