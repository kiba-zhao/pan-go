// Define stream for quic
package quic

import (
	"io"

	"github.com/quic-go/quic-go"
)

type stdQuicStream struct {
	quic.Stream
	conn   *stdQuicConn
	hangup bool
}

var _ = (io.ReadWriteCloser)((*stdQuicStream)(nil))

// Read implements the io.Reader interface for quic.Stream. It reads from the
// stream and calls OnStreamRead on the connection to update the stream's read
// offset.
func (s *stdQuicStream) Read(b []byte) (int, error) {
	n, err := s.Stream.Read(b)
	s.conn.OnStreamRead(s.StreamID(), n)
	return n, err
}

// Write implements the io.Writer interface for quic.Stream. It writes to the
// stream and calls OnStreamWrite on the connection to update the stream's write
// offset.
func (s *stdQuicStream) Write(b []byte) (int, error) {
	n, err := s.Stream.Write(b)
	s.conn.OnStreamWrite(s.StreamID(), n)
	return n, err
}

// Close closes the stdQuicStream and the underlying quic.Stream. If the stream is not
// flagged as hangup, it also informs the associated connection to close the stream.
// It returns any error encountered while closing the underlying quic.Stream.

func (s *stdQuicStream) Close() error {

	err := s.Stream.Close()
	if !s.hangup {
		s.conn.CloseStream(s)
	}

	return err
}
