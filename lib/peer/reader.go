package peer

import (
	"io"
)

type stdPeerReader struct {
	io.Reader
	haveRead bool
}

var _ = (io.Reader)((*stdPeerReader)(nil))

func (r *stdPeerReader) Read(p []byte) (n int, err error) {

	if !r.haveRead {
		r.haveRead = true
	}

	return r.Reader.Read(p)
}
