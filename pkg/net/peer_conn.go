package net

import (
	"context"
	"io"
)

type PeerStream interface {
	io.Reader
	io.Writer
	io.Closer
}

type PeerConn interface {
	PeerID() PeerID
	OpenStream(context.Context) (PeerStream, error)
	AcceptStream(context.Context) (PeerStream, error)
	Close() error
}
