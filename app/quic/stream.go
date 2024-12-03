package quic

import "github.com/quic-go/quic-go"

type quicStream struct {
	quic.Stream
	conn   *quicConn
	hangup bool
}

func (s *quicStream) Close() error {
	err := s.Stream.Close()
	if !s.hangup {
		s.conn.CloseStream(s)
	}

	return err
}
