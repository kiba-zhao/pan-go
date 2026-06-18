package net

import (
	"sync"

	"github.com/quic-go/quic-go"
)

type stdQuicStream struct {
	quic.Stream
	conn *stdQuicConn

	closed     bool
	closedLock sync.Mutex
}

func (s *stdQuicStream) Read(p []byte) (n int, err error) {
	n, err = s.Stream.Read(p)
	if err != nil {
		s.releaseForConn(false)
	}
	return n, err
}

func (s *stdQuicStream) Close() error {
	return s.releaseForConn(true)
}

func (s *stdQuicStream) releaseForConn(isReadable bool) error {
	s.closedLock.Lock()
	defer s.closedLock.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true

	if isReadable {
		s.Stream.CancelRead(quic.StreamErrorCode(quic.NoError))
	}
	err := s.Stream.Close()
	releaseOpenStreamForQuicConn(s.conn)
	return err
}
