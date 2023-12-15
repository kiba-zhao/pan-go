package peer

import (
	"context"
	"crypto/x509"
	"io"
)

type NodeServe interface {
	Accept(ctx context.Context) (Node, error)
	Close() error
}

type Node interface {
	Type() uint8
	Addr() []byte
	Certificate() *x509.Certificate
	AcceptNodeStream(ctx context.Context) (NodeStream, error)
	OpenNodeStream() (NodeStream, error)
	Close() error
}

type NodeDialer interface {
	Type() uint8
	Connect(addr []byte) (Node, error)
}

type NodeStreamCloser interface {
	io.Closer
	CloseRead() error
	CloseWrite() error
}

type NodeStream interface {
	io.Reader
	io.Writer
	NodeStreamCloser
}
