package node

import (
	"errors"
	"pan/lib/peer"
)

var ErrGuardAccessRefused = errors.New("guard.PeerGuard Error: Access Refused")

type PeerGuard struct {
	AppNodeService *AppNodeService
}

func (g *PeerGuard) Enabled() bool {
	return true
}

func (g *PeerGuard) Access(peerId peer.PeerID) error {
	return g.AppNodeService.AccessWithPeerID(peerId)
}
