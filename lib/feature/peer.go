package feature

import (
	"errors"
	"pan/lib/peer"
)

var ErrFeatureHelperPeerClusterNotFound = errors.New("FeatureHelper Error: PeerCluster Not Found")

type PeerTopic = peer.PeerAppModule

type PeerTopicProvider interface {
	PeerTopics() []PeerTopic
}

var _ = (peer.PeerAppModule)((*stdFeatureModule)(nil))

func (module *stdFeatureModule) SetupToPeer(app peer.PeerRouter) error {

	var errs []error
	featureHelper := module.featureHelper
	topics := getPeerTopics(featureHelper.feature)
	if len(topics) <= 0 {
		return nil
	}

	router := app.Route(featureHelper.PeerScope())
	for _, topic := range topics {
		err := topic.SetupToPeer(router)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func getPeerTopics(feature interface{}) []PeerTopic {
	if provider, ok := feature.(PeerTopicProvider); ok {
		return provider.PeerTopics()
	}
	return nil
}
