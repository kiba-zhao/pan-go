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

func (s *Settings) init() error {
	s.WebHost = "127.0.0.1"
	s.WebPort = 9002
	s.AppModule = "app"

	return nil
}

func (s *Settings) AppName() string {
	return "pan-go"
}

func (s *Settings) AppRoot() string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return path.Join(homePath, "."+s.AppName())
}
