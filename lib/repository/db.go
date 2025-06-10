package repository

import (
	"os"
	"pan/lib/config"
	"pan/lib/pkg"
	"path/filepath"
	"sync"

	"gorm.io/gorm"
)

type RepositoryDB = *gorm.DB

var dbPath string
var dbPathRW sync.RWMutex

func DBPath() string {
	dbPathRW.RLock()
	defer dbPathRW.RUnlock()
	return dbPath
}

func SetDBPath(path string) {
	dbPathRW.Lock()
	defer dbPathRW.Unlock()
	dbPath = path
}

func InitDBPath() error {
	cfgPath := config.RootPath()
	if len(cfgPath) > 0 {
		SetDBPath(cfgPath)
		return nil
	}

	err := config.InitRootPath()
	if err == nil {
		SetDBPath(config.RootPath())
	}
	return err
}

var tempdbPath string
var tempdbPathRW sync.RWMutex

func TempDBPath() string {
	tempdbPathRW.RLock()
	defer tempdbPathRW.RUnlock()
	return tempdbPath
}

func SetTempDBPath(path string) {
	tempdbPathRW.Lock()
	defer tempdbPathRW.Unlock()
	tempdbPath = path
}

func InitTempDBPath() error {
	tmpDir := filepath.Join(os.TempDir(), pkg.Name())
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return err
	}
	SetTempDBPath(tmpDir)
	return nil
}
