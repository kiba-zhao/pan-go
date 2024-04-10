package app

type AppSettings = *Settings

type Settings struct {
	RootPath string
	WebHost  string
	WebPort  int
}

func newDefaultSettings(rootPath string) AppSettings {
	settings := &Settings{}
	settings.RootPath = rootPath
	settings.WebHost = "127.0.0.1"
	settings.WebPort = 9002
	return settings
}
