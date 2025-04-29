package node

import (
	"iter"
	"pan/lib/peer"
)

type NetworkAddrGuide struct {
	NetworkAddrService *NetworkAddrService
}

func (guide *NetworkAddrGuide) LookupExplorerAddr(peerId peer.PeerID) (iter.Seq[string], error) {
	peerIdStr := peer.EncodePeerID(peerId)
	entitySeq, err := guide.NetworkAddrService.SeqByPeerId(peerIdStr)
	if err != nil || entitySeq == nil {
		return nil, err
	}
	return func(yield func(string) bool) {
		for entity := range entitySeq {
			if !yield(entity.Address) {
				break
			}
		}
	}, nil
}
