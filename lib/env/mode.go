package env

import (
	"os"
)

const (
	EnvMode     = "PANGO_MODE"
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)

func Mode() string {
	mode, ok := os.LookupEnv(EnvMode)
	if !ok {
		if os.Getenv("DEBUG") != "" {
			return DebugMode
		}
		return ReleaseMode
	}

	if mode != DebugMode && mode != ReleaseMode && mode != TestMode {
		return ReleaseMode
	}

	return mode
}
