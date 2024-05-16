package app

import (
	"os"
	"pan/runtime"
	"path"
)

func New() interface{} {

	settings := newDefaultSettings(getRootPath())
	config := NewConfig(settings, parseConfigPath)

	return runtime.NewModule(&runtime.Injector{}, config, NewNodeModule(), &webServer{})
}

func getRootPath() string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return path.Join(homePath, ".pan-go")
}

func parseConfigPath(settings AppSettings) string {
	return path.Join(settings.RootPath, "pan.toml")
}
