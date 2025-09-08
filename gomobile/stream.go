//go:build android || ios

package gomobile

import "io"

type GoMobileReadStream interface {
	io.Reader
}

type GoMobileStream interface {
	GoMobileReadStream
	io.Closer
}
