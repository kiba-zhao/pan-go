//go:build android || ios

package gomobile

import "pan/internal/log"

type Logger interface {
	log.Logger
}
