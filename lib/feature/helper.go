package feature

import "pan/lib/peer"

type stdFeatureHelper struct {
	PeerNetwork peer.PeerNetwork

	name    string
	feature interface{}
}

func (helper *stdFeatureHelper) FeatureName() string {
	return helper.name
}

func (helper *stdFeatureHelper) WebScope() string {
	return "/api/" + helper.FeatureName()
}

func (helper *stdFeatureHelper) PeerScope() []byte {
	return []byte(helper.FeatureName() + ".")
}

func (helper *stdFeatureHelper) SerlvetScope() []byte {
	return helper.PeerScope()
}
