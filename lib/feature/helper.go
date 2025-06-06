package feature

type stdFeatureHelper struct {
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
