package gomobile

import "io"

type GoMobileReadStream = io.Reader
type GoMobileStream interface {
	GoMobileReadStream
	io.Closer
}
