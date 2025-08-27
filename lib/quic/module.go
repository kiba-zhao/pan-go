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
	PeerGuard   peer.PeerGuard
	PeerNetwork peer.PeerNetwork

	configurer QuicConfigurer
	network    *stdQuicNetwork
	explorer   *stdQuicExplorer
	agent      *stdQuicAgent
	server     *stdQuicServer
}

func New() interface{} {

	network := &stdQuicNetwork{}
	network.connMgr = &stdQuicConnMgr{}
	network.routeMgr = &stdQuicRouteMgr{}

	agent := &stdQuicAgent{}
	agent.network = network
	network.agent = agent

	explorer := &stdQuicExplorer{}
	explorer.network = network

	server := &stdQuicServer{}
	server.network = network
	server.reloadChan = make(chan struct{}, 1)

	provider := &stdQuicProvider{}
	provider.transportMap = make(map[string]*QuicTransport)
	network.provider = provider
	server.provider = provider

	logger := log.Default()
	network.logger = logger
	agent.logger = logger
	server.logger = logger

	configurer := config.NewConfigurer[QuicConfig](logger)

	module := &stdQuicModule{}
	module.agent = agent
	module.network = network
	module.explorer = explorer
	module.server = server
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
	qm.server.Setup(config)
}

var _ = (bootstrap.ReadyModule)((*stdQuicModule)(nil))

func (qm *stdQuicModule) Ready(ctx context.Context) error {
	qm.network.SetupPeerNetwork(qm.PeerNetwork)
	qm.network.SetupPeerGuard(qm.PeerGuard)

	qm.PeerNetwork.RegisterPeerTransport(qm.network)
	defer qm.PeerNetwork.UnregisterPeerTransport(qm.network)
	qm.PeerNetwork.RegisterPeerTransport(qm.explorer)
	defer qm.PeerNetwork.UnregisterPeerTransport(qm.explorer)

	qm.configurer.Subscribe(qm)
	defer qm.configurer.Unsubscribe(qm)

	return qm.server.ListenAndServe(ctx)
}
