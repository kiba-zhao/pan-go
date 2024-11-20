package peer

import (
	"context"
	"io"
)

type PeerID = []byte
type PeerType = uint8
type PeerResourceID = []byte

const (
	PeerTypeAlive PeerType = iota + 1
	PeerTypeReachable
)

type PeerNode interface {
	PeerID() PeerID
	PeerType() PeerType
	Do(context.Context, io.Reader) (io.ReadCloser, error)
	Close() error
	PeerResourceID() PeerResourceID
	IsIdle() bool
}
