package net

import "github.com/quic-go/quic-go"

type stdQuicStream struct {
	quic.Stream
}

func (s *stdQuicStream) Close() error {
	s.Stream.CancelRead(quic.StreamErrorCode(quic.NoError))
	return s.Stream.Close()
}
