package net

type PeerAppModule interface {
	SetupToPeer(PeerServlet) error
}

type PeerTopic interface {
	SetupToPeer(PeerServletRouter) error
}

type PeerTopicProvider interface {
	PeerTopics() []PeerTopic
}

type PeerRouteModule interface {
	PeerRouteScope() PeerServletScope
}
