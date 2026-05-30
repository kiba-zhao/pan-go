package user

import (
	"io"
	"iter"
	"slices"
)

type IterStream struct {
	buffer []byte
	next   func() ([]byte, error, bool)
	stop   func()
	err    error
}

func NewIterStreamForSeq2(seq iter.Seq2[[]byte, error]) *IterStream {
	var stream IterStream
	next, stop := iter.Pull2(seq)
	stream.next = next
	stream.stop = stop
	return &stream
}

var _ = (io.ReadCloser)((*IterStream)(nil))

func (stream *IterStream) Close() error {
	stream.stop()
	stream.err = io.ErrClosedPipe
	return nil
}

func (stream *IterStream) Read(p []byte) (n int, err error) {
	maxLen := len(p)
	bufferLen := len(stream.buffer)
	err = stream.err
	if err == nil && maxLen > bufferLen {
		nextBytes, err, done := stream.next()
		if err != nil {
			stream.err = err
		} else if len(nextBytes) > 0 {
			stream.buffer = slices.Concat(stream.buffer, nextBytes)
			bufferLen = len(stream.buffer)
		}
		if done {
			err = io.EOF
			stream.err = err
		}
	}

	if err != nil && bufferLen == 0 {
		return 0, err
	}

	n = min(maxLen, bufferLen)
	copy(p, stream.buffer[:n])
	if n == bufferLen {
		stream.buffer = nil
	} else {
		stream.buffer = slices.Clone(stream.buffer[n:])
	}

	return n, nil
}
