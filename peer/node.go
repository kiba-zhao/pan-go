package peer

import (
	"context"
	"crypto/x509"
	"io"
)

type NodeType = uint8

type NodeTypeComponent interface {
	Type() NodeType
}

type NodeServe interface {
	Accept(ctx context.Context) (Node, error)
	Close() error
}

type Node interface {
	NodeTypeComponent
	Addr() []byte
	Certificate() *x509.Certificate
	AcceptNodeStream(ctx context.Context) (NodeStream, error)
	OpenNodeStream() (NodeStream, error)
	Close() error
}

type NodeDialer interface {
	NodeTypeComponent
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

type NodeHandshake interface {
	NodeTypeComponent
	Handshake() []byte
}
