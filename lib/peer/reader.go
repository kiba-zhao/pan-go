package peer

import (
	"io"
)

type peerReader struct {
	io.Reader
	haveRead bool
}

func (r *peerReader) Read(p []byte) (n int, err error) {

	if !r.haveRead {
		r.haveRead = true
	}

	return r.Reader.Read(p)
}
