package app

import "path"

type AppSettings = *Settings

type Settings struct {
	RootPath         string
	WebAddress       []string
	NodeAddress      []string
	PrivateKeyPath   string
	CertificatePath  string
	BroadcastAddress []string
	PublicAddress    []string
}

func newDefaultSettings(rootPath string) AppSettings {
	settings := &Settings{}
	settings.RootPath = rootPath
	settings.WebAddress = []string{"127.0.0.1:9002"}
	settings.NodeAddress = []string{"0.0.0.0:9000"}
	settings.BroadcastAddress = []string{"224.0.0.120:9100"}
	settings.PrivateKeyPath = path.Join(rootPath, "key.pem")
	settings.CertificatePath = path.Join(rootPath, "cert.pem")
	settings.PublicAddress = settings.NodeAddress

	return settings
}
