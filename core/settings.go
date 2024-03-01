package core

import (
	"os"
	"path"
)

type Settings struct {
	WebHost   string
	WebPort   int
	AppModule string
}

func (settings *Settings) init() error {
	settings.WebHost = "127.0.0.1"
	settings.WebPort = 9002
	settings.AppModule = "app"
	return nil
}

func (settings *Settings) AppName() string {
	return "pan-go"
}

func (settings *Settings) AppRoot() string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return path.Join(homePath, "."+settings.AppName())
}
