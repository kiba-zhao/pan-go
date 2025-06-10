package web

import "sync/atomic"

var webMode atomic.Bool

func IsWebMode() bool {
	return webMode.Load()
}

func SetWebModeEnabled() {
	webMode.Store(true)
}

func SetWebModeDisabled() {
	webMode.Store(false)
}
