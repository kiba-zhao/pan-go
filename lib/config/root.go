package config

import (
	"os"
	"pan/lib/pkg"
	"path/filepath"
	"sync"
)

var cfgPath string
var cfgPathRW sync.RWMutex

func RootPath() string {
	cfgPathRW.RLock()
	defer cfgPathRW.RUnlock()

	return cfgPath
}

func SetRootPath(path string) {
	cfgPathRW.Lock()
	defer cfgPathRW.Unlock()
	cfgPath = path
}

func InitRootPath() error {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	SetRootPath(filepath.Join(homePath, "."+pkg.Name()))
	return nil
}
