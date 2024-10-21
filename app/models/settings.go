package models

import "pan/app/config"

type Settings struct {
	config.Settings
	NodeID   string `json:"nodeId" form:"nodeId"`
	RootPath string `json:"rootPath" form:"rootPath"`
}

type SettingsFields struct {
	Name             string   `form:"name" json:"name"  binding:"omitempty"`
	WebAddress       []string `form:"webAddress" json:"webAddress"  binding:"omitempty"`
	NodeAddress      []string `form:"nodeAddress" json:"nodeAddress"  binding:"omitempty"`
	BroadcastAddress []string `form:"broadcastAddress" json:"broadcastAddress"  binding:"omitempty"`
	PublicAddress    []string `form:"publicAddress" json:"publicAddress"  binding:"omitempty"`
	GuardEnabled     *bool    `form:"guardEnabled" json:"guardEnabled"  binding:"omitempty"`
	GuardAccess      *bool    `form:"guardAccess" json:"guardAccess"  binding:"omitempty"`
}
