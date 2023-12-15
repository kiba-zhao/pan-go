package peer_test

import (
	"io"
)

type TestCloseFn func() error

type TestNodeStream struct {
	io.Closer
	io.Reader
	io.Writer
	CloseReadFn  TestCloseFn
	CloseWriteFn TestCloseFn
}

// CloseRead ...
func (ns *TestNodeStream) CloseRead() error {
	return ns.CloseReadFn()
}

// CloseWrite ...
func (ns *TestNodeStream) CloseWrite() error {
	return ns.CloseWriteFn()
}
