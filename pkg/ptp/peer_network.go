package ptp

import (
	"context"
)

type PeerID = []byte

type PeerNetwork interface {
	RoundTrip(context.Context, PeerID) (PeerStream, error)
	Connect(context.Context, PeerID) (PeerConn, error)
	Reuse(PeerConn) error
}
