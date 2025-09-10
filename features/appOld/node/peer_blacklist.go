package node

import (
	"errors"
	"pan/lib/peer"
)

var ErrGuardAccessRefused = errors.New("guard.PeerGuard Error: Access Refused")

type PeerBlackList struct {
	AppNodeService *AppNodeService
}

func (g *PeerBlackList) Enabled() bool {
	return true
}

func (g *PeerBlackList) Has(peerId peer.PeerID) bool {
	return g.AppNodeService.AccessWithPeerID(peerId) == nil
}
