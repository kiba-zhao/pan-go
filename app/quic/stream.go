package quic

import "github.com/quic-go/quic-go"

type quicPeerStream struct {
	quic.Stream
	quicPeerNode *quicPeerNode
	closed       bool
	isServe      bool
}

func (qs *quicPeerStream) Read(p []byte) (n int, err error) {
	n, err = qs.Stream.Read(p)
	if err != nil && !qs.isServe {
		qs.Close()
	}
	return
}

func (qs *quicPeerStream) Close() error {
	if qs.closed {
		return nil
	}
	qs.closed = true
	qs.quicPeerNode.decreaseStream()
	var err error
	if qs.isServe {
		err = qs.Stream.Close()
	} else {
		qs.CancelRead(quic.StreamErrorCode(quic.NoError))
	}

	return err
}
