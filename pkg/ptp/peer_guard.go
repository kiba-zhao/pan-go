package ptp

import (
	"slices"
	"sync"
)

const (
	PeerGuardDefault = uint8(iota + 1)
	PeerGuardAllow
	PeerGuardDeny
)

type PeerGuardGuide interface {
	IsAccessible(PeerID) bool
	IsDenied(PeerID) bool
}

type stdPeerGuard struct {
	guides []PeerGuardGuide
	rw     sync.RWMutex
}

func (guard *stdPeerGuard) setup(guides []PeerGuardGuide) {
	guard.rw.Lock()
	defer guard.rw.Unlock()
	guard.guides = guides
}

func (guard *stdPeerGuard) check(peerId PeerID) uint8 {
	var guides []PeerGuardGuide
	guard.rw.RLock()
	if len(guard.guides) > 0 {
		guides = slices.Clone(guard.guides)
	} else {
		guard.rw.RUnlock()
		return PeerGuardDefault
	}
	guard.rw.RUnlock()

	isAllow := false
	isDenied := false
	for _, guide := range guides {
		isDenied = guide.IsDenied(peerId)
		if isDenied {
			break
		}
		if !isAllow {
			isAllow = guide.IsAccessible(peerId)
		}
	}

	if isDenied {
		return PeerGuardDeny
	}
	if isAllow {
		return PeerGuardAllow
	}
	return PeerGuardDefault
}
