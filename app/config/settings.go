package config

import (
	"os"
	"path"
)

type AppSettings = *Settings

type Settings struct {
	RootPath         string   `json:"rootPath" form:"rootPath"`
	Name             string   `json:"name" form:"name"`
	WebAddress       []string `json:"webAddress" form:"webAddress"`
	NodeAddress      []string `json:"nodeAddress" form:"nodeAddress"`
	BroadcastAddress []string `json:"broadcastAddress" form:"broadcastAddress"`
	PublicAddress    []string `json:"publicAddress" form:"publicAddress"`
}

func newDefaultSettings() AppSettings {
	settings := &Settings{}
	settings.RootPath = generateRootPath()
	settings.Name = generateName()
	settings.WebAddress = []string{"127.0.0.1:9002"}
	settings.NodeAddress = []string{"0.0.0.0:9000"}
	settings.BroadcastAddress = []string{"224.0.0.120:9100"}
	settings.PublicAddress = settings.NodeAddress

	return settings
}

func generateRootPath() string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return path.Join(homePath, ".pan-go")
}

func generateName() string {
	name, err := os.Hostname()
	if err == nil {
		return name
	}

	return "pan-go"
}
