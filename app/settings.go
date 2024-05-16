package app

import "path"

type AppSettings = *Settings

type Settings struct {
	RootPath        string
	WebHost         string
	WebPort         int
	Host            string
	Port            int
	PrivateKeyPath  string
	CertificatePath string
}

func newDefaultSettings(rootPath string) AppSettings {
	settings := &Settings{}
	settings.RootPath = rootPath
	settings.WebHost = "127.0.0.1"
	settings.WebPort = 9002
	settings.Host = "0.0.0.0"
	settings.Port = 9000
	settings.PrivateKeyPath = path.Join(rootPath, "key.pem")
	settings.CertificatePath = path.Join(rootPath, "cert.pem")
	return settings
}
