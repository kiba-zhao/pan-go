package feature

import (
	"errors"
	"pan/lib/peer"
)

var ErrFeatureHelperPeerNetworkNotFound = errors.New("FeatureHelper Error: Peer Network Not Found")

type PeerTopic = peer.PeerAppModule

type PeerTopicProvider interface {
	PeerTopics() []PeerTopic
}

func getPeerTopics(feature interface{}) []PeerTopic {
	if provider, ok := feature.(PeerTopicProvider); ok {
		return provider.PeerTopics()
	}
	return nil
}
