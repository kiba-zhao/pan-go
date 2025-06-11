package config

import (
	"errors"
	"os"
	"sync"
)

var hostname string
var hostnameRW sync.RWMutex
var hostnameAlready bool
var errHostNameConflict = errors.New("config.HostName Error: Conflict")

func HostName() string {
	hostnameRW.RLock()
	defer hostnameRW.RUnlock()
	return hostname
}

func InitHostName(name string) error {
	hostnameRW.Lock()
	defer hostnameRW.Unlock()
	if hostnameAlready {
		return errHostNameConflict
	}
	hostname = name
	hostnameAlready = true
	return nil
}

func initHostNameAsDefault() error {
	name, err := os.Hostname()
	if err == nil {
		return InitHostName(name)
	}
	return err
}
