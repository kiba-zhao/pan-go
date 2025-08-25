package peer

import (
	"pan/lib/log"
	"slices"
	"sync"
)

type PeerBlackList interface {
	Enabled() bool
	Has(PeerID) bool
}

type PeerGuard interface {
	Access(PeerID) bool
	RegisterBlackList(PeerBlackList)
	UnregisterBlackList(PeerBlackList)
}

type stdPeerGuard struct {
	logger log.Logger

	blackLists   []PeerBlackList
	blackListsRW sync.RWMutex
}

var _ = (PeerGuard)((*stdPeerGuard)(nil))

func (guard *stdPeerGuard) Access(peerId PeerID) bool {
	guard.blackListsRW.RLock()
	defer guard.blackListsRW.RUnlock()

	blackLists := guard.blackLists
	if len(blackLists) <= 0 {
		return true
	}

	for _, blackList := range blackLists {
		if !blackList.Enabled() {
			continue
		}
		if blackList.Has(peerId) {
			return false
		}
	}

	return true
}

func (guard *stdPeerGuard) RegisterBlackList(blackList PeerBlackList) {
	guard.blackListsRW.Lock()
	defer guard.blackListsRW.Unlock()
	guard.blackLists = append(guard.blackLists, blackList)
}

func (guard *stdPeerGuard) UnregisterBlackList(blackList PeerBlackList) {
	guard.blackListsRW.Lock()
	defer guard.blackListsRW.Unlock()
	guard.blackLists = slices.DeleteFunc(guard.blackLists, func(b PeerBlackList) bool {
		return b == blackList
	})
}
