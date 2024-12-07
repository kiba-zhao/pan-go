package quic

import (
	"github.com/quic-go/quic-go"
)

type quicStream struct {
	quic.Stream
	conn   *quicConn
	hangup bool
}

func (s *quicStream) Read(b []byte) (int, error) {
	n, err := s.Stream.Read(b)
	s.conn.OnStreamRead(s.StreamID(), n)
	return n, err
}

func (s *quicStream) Write(b []byte) (int, error) {
	n, err := s.Stream.Write(b)
	s.conn.OnStreamWrite(s.StreamID(), n)
	return n, err
}

func (s *quicStream) Close() error {

	err := s.Stream.Close()
	if !s.hangup {
		s.conn.CloseStream(s)
	}

	return err
}
