package appsettings

import "pan/app/config"

type AppSettings struct {
	config.Settings
	PeerID   string `json:"peerId" form:"peerId"`
	RootPath string `json:"rootPath" form:"rootPath"`
}

type AppSettingsFields struct {
	Name             string   `form:"name" json:"name"  binding:"omitempty"`
	WebAddress       []string `form:"webAddress" json:"webAddress"  binding:"omitempty"`
	PeerAddress      []string `form:"peerAddress" json:"peerAddress"  binding:"omitempty"`
	BroadcastAddress []string `form:"broadcastAddress" json:"broadcastAddress"  binding:"omitempty"`
	PublicAddress    []string `form:"publicAddress" json:"publicAddress"  binding:"omitempty"`
	GuardEnabled     *bool    `form:"guardEnabled" json:"guardEnabled"  binding:"omitempty"`
	GuardAccess      *bool    `form:"guardAccess" json:"guardAccess"  binding:"omitempty"`
}
