package config

import (
	"errors"
	"os"
	"pan/lib/pkg"
	"path/filepath"
	"sync"
)

var cfgPath string
var cfgPathRW sync.RWMutex
var cfgPathAlready bool
var errRootPathConflict = errors.New("config.RootPath Error: Conflict")

func RootPath() string {
	cfgPathRW.RLock()
	defer cfgPathRW.RUnlock()

	return cfgPath
}

func InitRootPath(path string) error {
	cfgPathRW.Lock()
	defer cfgPathRW.Unlock()
	if cfgPathAlready {
		return errRootPathConflict
	}
	cfgPath = path
	cfgPathAlready = true
	return nil
}

func initRootPathAsDefault() error {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	InitRootPath(filepath.Join(homePath, "."+pkg.Name()))
	return nil
}
