// Define stream for quic
package quic

import (
	"github.com/quic-go/quic-go"
)

type quicStream struct {
	quic.Stream
	conn   *quicConn
	hangup bool
}

// Read implements the io.Reader interface for quic.Stream. It reads from the
// stream and calls OnStreamRead on the connection to update the stream's read
// offset.
func (s *quicStream) Read(b []byte) (int, error) {
	n, err := s.Stream.Read(b)
	s.conn.OnStreamRead(s.StreamID(), n)
	return n, err
}

// Write implements the io.Writer interface for quic.Stream. It writes to the
// stream and calls OnStreamWrite on the connection to update the stream's write
// offset.
func (s *quicStream) Write(b []byte) (int, error) {
	n, err := s.Stream.Write(b)
	s.conn.OnStreamWrite(s.StreamID(), n)
	return n, err
}

// Close closes the quicStream and the underlying quic.Stream. If the stream is not
// flagged as hangup, it also informs the associated connection to close the stream.
// It returns any error encountered while closing the underlying quic.Stream.

func (s *quicStream) Close() error {

	err := s.Stream.Close()
	if !s.hangup {
		s.conn.CloseStream(s)
	}

	return err
}
