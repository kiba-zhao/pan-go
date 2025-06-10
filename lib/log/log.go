package log

import (
	"sync"
)

var logger Logger
var loggerOnce sync.Once

func Default() Logger {
	return logger
}

func InitDefault(defaultLogger Logger) {
	loggerOnce.Do(func() {
		logger = defaultLogger
	})
}
