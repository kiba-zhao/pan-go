package net

import (
	"context"
	"io"
)

type PeerConnState interface {
	Session() ([]byte, error)
}

type PeerStream interface {
	io.Reader
	io.Writer
	io.Closer
}

type PeerConn interface {
	ConnState() PeerConnState
	OpenStream(context.Context) (PeerStream, error)
	AcceptStream(context.Context) (PeerStream, error)
}
