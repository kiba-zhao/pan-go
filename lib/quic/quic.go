// Package quic provides the quic module
package quic

func New() interface{} {

	var module quicPeerModule
	module.connMgr = &quicConnMgr{}
	module.routeMgr = &quicRouteMgr{}

	var agent quicPeerAgent
	module.agent = &agent
	agent.quicPeerModule = &module

	return &module
}
