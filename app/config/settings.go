package config

import (
	"os"
)

type AppSettings = *Settings

type Settings struct {
	Name             string   `json:"name" form:"name"`
	WebAddress       []string `json:"webAddress" form:"webAddress"`
	NodeAddress      []string `json:"nodeAddress" form:"nodeAddress"`
	BroadcastAddress []string `json:"broadcastAddress" form:"broadcastAddress"`
	PublicAddress    []string `json:"publicAddress" form:"publicAddress"`
	GuardEnabled     bool     `json:"guardEnabled" form:"guardEnabled"`
	GuardAccess      bool     `json:"guardAccess" form:"guardAccess"`
}

func newDefaultSettings() AppSettings {
	settings := &Settings{}
	settings.Name = generateName()
	settings.WebAddress = []string{"127.0.0.1:9002"}
	settings.NodeAddress = []string{"0.0.0.0:9000"}
	settings.BroadcastAddress = []string{"224.0.0.120:9100"}
	settings.PublicAddress = settings.NodeAddress
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
