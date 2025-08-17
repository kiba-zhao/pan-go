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
	AppConfig   config.AppConfig
	PeerConfig  peer.PeerConfig
	PeerCluster peer.PeerCluster

	cluster  *stdQuicCluster
	explorer *stdQuicExplorer
	agent    *stdQuicAgent
	server   *stdQuicServer
}

func New() interface{} {

	logger := log.Default()

	var module stdQuicModule
	cluster := &stdQuicCluster{}
	module.cluster = cluster
	cluster.connMgr = &stdQuicConnMgr{}
	cluster.routeMgr = &stdQuicRouteMgr{}

	explorer := &stdQuicExplorer{}
	module.explorer = explorer
	explorer.cluster = cluster

	agent := &stdQuicAgent{}
	module.agent = agent
	agent.logger = logger
	agent.cluster = cluster
	cluster.agent = agent

	server := &stdQuicServer{}
	module.server = server
	server.logger = logger
	server.cluster = cluster
	cluster.server = server
	server.reloadChan = make(chan struct{}, 1)

	return &module
}

var _ = (injection.ComponentProvider)((*stdQuicModule)(nil))

func (qm *stdQuicModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(qm, injection.ComponentNoneScope),
		injection.NewComponent[QuicCluster](qm.cluster, injection.ComponentExternalScope),
		injection.NewComponent[QuicExplorer](qm.explorer, injection.ComponentExternalScope),
	}
}

var _ = (config.AppConfigListener)((*stdQuicModule)(nil))

func (qm *stdQuicModule) OnConfigUpdated(settings config.AppSettings) {
	if settings == nil {
		return
	}

	qm.server.SetPort(settings.PeerPort)

}

var _ = (peer.PeerConfigListener)((*stdQuicModule)(nil))

func (qm *stdQuicModule) OnPeerConfigUpdated(settings *peer.PeerSettings) {
	if settings == nil {
		return
	}

	certificate := settings.Certificate()
	qm.cluster.SetCertificate(certificate)
	qm.server.SetCertificate(certificate)
}

var _ = (bootstrap.ReadyModule)((*stdQuicModule)(nil))

func (qm *stdQuicModule) Ready(ctx context.Context) error {
	qm.cluster.SetPeerCluster(qm.PeerCluster)

	qm.cluster.peerCluster.RegisterPeerNetwork(qm.cluster)
	defer qm.cluster.peerCluster.UnregisterPeerNetwork(qm.cluster)
	qm.cluster.peerCluster.RegisterPeerNetwork(qm.explorer)
	defer qm.cluster.peerCluster.UnregisterPeerNetwork(qm.explorer)

	qm.AppConfig.Subscribe(qm)
	defer qm.AppConfig.Unsubscribe(qm)

	qm.PeerConfig.Subscribe(qm)
	defer qm.PeerConfig.Unsubscribe(qm)

	return qm.server.ListenAndServe(ctx)
}
