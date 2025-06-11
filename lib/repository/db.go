package repository

import (
	"errors"
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
var dbPathAlready bool
var errDBPathConflict = errors.New("repository.DBPath Error: Conflict")
var errRootPathUnavailable = errors.New("repository.DBPath Error: Root Path Unavailable")

func DBPath() string {
	dbPathRW.RLock()
	defer dbPathRW.RUnlock()
	return dbPath
}

func InitDBPath(path string) error {
	dbPathRW.Lock()
	defer dbPathRW.Unlock()
	if dbPathAlready {
		return errDBPathConflict
	}
	dbPath = path
	dbPathAlready = true
	return nil
}

func initDBPathAsDefault() error {
	cfgPath := config.RootPath()
	if len(cfgPath) > 0 {
		return InitDBPath(cfgPath)
	}

	return errRootPathUnavailable
}

var tempdbPath string
var tempdbPathRW sync.RWMutex
var tempdbPathAlready bool
var errTempDBPathConflict = errors.New("repository.TempDBPath Error: Conflict")

func TempDBPath() string {
	tempdbPathRW.RLock()
	defer tempdbPathRW.RUnlock()
	return tempdbPath
}

func InitTempDBPath(path string) error {
	tempdbPathRW.Lock()
	defer tempdbPathRW.Unlock()
	if tempdbPathAlready {
		return errTempDBPathConflict
	}
	tempdbPath = path
	tempdbPathAlready = true
	return nil
}

func initTempDBPathAsDefault() error {
	tmpDir := filepath.Join(os.TempDir(), pkg.Name())
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return err
	}
	return InitTempDBPath(tmpDir)
}
