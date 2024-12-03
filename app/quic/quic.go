package quic

func New() QuicPeerModule {

	var module quicPeerModule
	module.connMgr = &quicConnMgr{}
	module.routeMgr = &quicRouteMgr{}

	var broadcast quicPeerBroadcast
	broadcast.quicPeerModule = &module
	module.quicPeerBroadcast = &broadcast

	return &module
}
