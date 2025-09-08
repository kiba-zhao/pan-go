//go:build android || ios

package gomobile

import "pan/lib/log"

type Logger interface {
	log.Logger
}
