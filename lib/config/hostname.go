package config

import (
	"os"
	"pan/lib/pkg"
	"sync"
)

var hostname string
var hostnameRW sync.RWMutex

func HostName() string {
	hostnameRW.RLock()
	defer hostnameRW.RUnlock()
	return hostname
}

func SetHostName(name string) {
	hostnameRW.Lock()
	defer hostnameRW.Unlock()
	hostname = name
}

func InitHostName() error {
	name, err := os.Hostname()
	if err == nil {
		SetHostName(name)
	} else {
		SetHostName(pkg.Name())
	}

	return err
}
