package config

import (
	"os"
	"path"
)

type AppSettings = *Settings

type Settings struct {
	RootPath         string
	Name             string
	WebAddress       []string
	NodeAddress      []string
	PrivateKeyPath   string
	CertificatePath  string
	BroadcastAddress []string
	PublicAddress    []string
}

func newDefaultSettings() AppSettings {
	settings := &Settings{}
	settings.RootPath = generateRootPath()
	settings.Name = generateName()
	settings.WebAddress = []string{"127.0.0.1:9002"}
	settings.NodeAddress = []string{"0.0.0.0:9000"}
	settings.BroadcastAddress = []string{"224.0.0.120:9100"}
	settings.PrivateKeyPath = path.Join(settings.RootPath, "key.pem")
	settings.CertificatePath = path.Join(settings.RootPath, "cert.pem")
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
