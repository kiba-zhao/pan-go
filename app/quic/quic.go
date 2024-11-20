package quic

func New() QuicPeerModule {

	var broadcast quicPeerBroadcast
	var module quicPeerModule
	module.quicPeerBroadcast = &broadcast
	broadcast.quicPeerModule = &module

	return &module
}
