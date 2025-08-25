package quic

import (
	"context"
	"pan/lib/bootstrap"
	"pan/lib/config"
	"pan/lib/injection"
	"pan/lib/log"
	"pan/lib/peer"
)

type stdQuicModule struct {
	PeerGuard  peer.PeerGuard
	PeerServer peer.PeerServer
	PeerClient peer.PeerClient

	configurer QuicConfigurer
	network    *stdQuicNetwork
	explorer   *stdQuicExplorer
	agent      *stdQuicAgent
	server     *stdQuicServer
}

func New() interface{} {

	logger := log.Default()

	network := &stdQuicNetwork{}
	network.logger = logger
	network.connMgr = &stdQuicConnMgr{}
	network.routeMgr = &stdQuicRouteMgr{}

	agent := &stdQuicAgent{}
	agent.logger = logger
	agent.network = network
	network.agent = agent

	explorer := &stdQuicExplorer{}
	explorer.network = network

	server := &stdQuicServer{}
	server.logger = logger
	server.network = network
	server.reloadChan = make(chan struct{}, 1)

	provider := &stdQuicTransportProvider{}
	provider.transportMap = make(map[string]*QuicTransport)
	network.provider = provider
	server.provider = provider

	module := &stdQuicModule{}
	module.agent = agent
	module.network = network
	module.explorer = explorer
	module.server = server

	configurer := config.NewConfigurer[QuicConfig](logger)
	module.configurer = configurer

	return module
}

var _ = (injection.ComponentProvider)((*stdQuicModule)(nil))

func (qm *stdQuicModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(qm, injection.ComponentNoneScope),
		injection.NewComponent[QuicNetwork](qm.network, injection.ComponentExternalScope),
		injection.NewComponent[QuicExplorer](qm.explorer, injection.ComponentExternalScope),
		injection.NewComponent[QuicConfigurer](qm.configurer, injection.ComponentExternalScope),
	}
}

var _ = (QuicConfigListener)((*stdQuicModule)(nil))

func (qm *stdQuicModule) OnConfigUpdated(config QuicConfig) {
	if config == nil {
		return
	}
	qm.network.Setup(config)
	qm.server.Reload(config)
}

var _ = (bootstrap.ReadyModule)((*stdQuicModule)(nil))

func (qm *stdQuicModule) Ready(ctx context.Context) error {
	qm.network.SetPeerServer(qm.PeerServer)
	qm.network.SetPeerGuard(qm.PeerGuard)

	qm.PeerClient.RegisterPeerTransport(qm.network)
	defer qm.PeerClient.UnregisterPeerTransport(qm.network)
	qm.PeerClient.RegisterPeerTransport(qm.explorer)
	defer qm.PeerClient.UnregisterPeerTransport(qm.explorer)

	qm.configurer.Subscribe(qm)
	defer qm.configurer.Unsubscribe(qm)

	return qm.server.ListenAndServe(ctx)
}
