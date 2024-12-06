package quic

import "pan/runtime"

func New() interface{} {

	var module quicPeerModule
	module.connMgr = &quicConnMgr{}
	module.routeMgr = &quicRouteMgr{}

	var broadcast quicPeerBroadcast
	broadcast.quicPeerModule = &module
	module.quicPeerBroadcast = &broadcast

	var agent quicPeerAgent
	module.agent = &agent
	agent.quicPeerModule = &module

	return runtime.NewModule(&module, &broadcast)
}
