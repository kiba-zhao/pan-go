// Package quic provides the quic module
package quic

import (
	"pan/lib/runtime"
)

func New() interface{} {

	var module quicPeerModule
	module.connMgr = &quicConnMgr{}
	module.routeMgr = &quicRouteMgr{}

	var agent quicPeerAgent
	module.agent = &agent
	agent.quicPeerModule = &module

	var explorer QuicExplorer
	explorer.quicModule = &module

	return runtime.NewModule(&module, &explorer)
}
