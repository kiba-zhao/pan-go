package broadcast

import (
	"context"
	"net"
	"pan/lib/bootstrap"
	"pan/lib/config"
	"pan/lib/injection"
	"pan/lib/log"
	"pan/lib/peer"
	"pan/lib/quic"
	"sync"
)

type BroadcastModule interface {
	SetStore(store BroadcastStore)
}

type stdBroadcastModule struct {
	AppConfig   config.AppConfig
	PeerConfig  peer.PeerConfig
	QuicCluster quic.QuicCluster

	cluster *stdBroadcastCluster
	agent   *stdBroadcastAgent
	server  *stdBroadcastServer

	store   BroadcastStore
	storeRW sync.RWMutex
}

func New() interface{} {
	module := &stdBroadcastModule{}

	cluster := &stdBroadcastCluster{}
	module.cluster = cluster
	cluster.SetDeliverLimitSize(65535)
	setDefualtMTU(cluster)

	agent := &stdBroadcastAgent{}
	module.agent = agent
	agent.reloadChan = make(chan struct{}, 1)
	agent.cluster = cluster
	cluster.RegisterServeModule(agent)

	server := &stdBroadcastServer{}
	module.server = server
	server.reloadChan = make(chan struct{}, 1)
	server.cluster = cluster

	logger := log.Default()
	cluster.logger = logger
	agent.logger = logger
	server.logger = logger

	return module
}

var _ = (BroadcastModule)((*stdBroadcastModule)(nil))

func (module *stdBroadcastModule) SetStore(store BroadcastStore) {
	module.storeRW.Lock()
	defer module.storeRW.Unlock()
	module.store = store
}

func (module *stdBroadcastModule) Store() BroadcastStore {
	module.storeRW.RLock()
	defer module.storeRW.RUnlock()
	return module.store
}

var _ = (injection.ComponentProvider)((*stdBroadcastModule)(nil))

func (module *stdBroadcastModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(module, injection.ComponentNoneScope),
		injection.NewComponent[BroadcastModule](module, injection.ComponentExternalScope),
		injection.NewComponent[BroadcastCluster](module.cluster, injection.ComponentExternalScope),
	}
}

var _ = (bootstrap.ReadyModule)((*stdBroadcastModule)(nil))

func (module *stdBroadcastModule) Ready(ctx context.Context) error {
	module.agent.SetQuicCluster(module.QuicCluster)

	module.AppConfig.Subscribe(module)
	defer module.AppConfig.Unsubscribe(module)

	module.PeerConfig.Subscribe(module)
	defer module.PeerConfig.Unsubscribe(module)

	var wg sync.WaitGroup
	wg.Add(2)
	go func(server *stdBroadcastServer) {
		defer wg.Done()
		server.ListenAndServe(ctx)
	}(module.server)

	go func(agent *stdBroadcastAgent) {
		defer wg.Done()
		agent.Deliver(ctx)
	}(module.agent)

	wg.Wait()
	return nil
}

var _ = (config.AppConfigListener)((*stdBroadcastModule)(nil))

func (module *stdBroadcastModule) OnConfigUpdated(settings config.AppSettings) {
	module.cluster.SetDeliverAddr(settings.BroadcastAddress)
	module.agent.SetPublicAddrs(settings.PublicAddress)
	module.server.SetAddrs(settings.BroadcastAddress)
}

var _ = (peer.PeerConfigListener)((*stdBroadcastModule)(nil))

func (module *stdBroadcastModule) OnPeerConfigUpdated(settings *peer.PeerSettings) {
	module.agent.SetPeerSettings(settings)
}

func setDefualtMTU(cluster *stdBroadcastCluster) {
	ifs, err := net.Interfaces()

	if err != nil {
		cluster.SetMTU(1500)
		return
	}

	mtu := cluster.DeliverLimitSize()
	for _, i := range ifs {
		if addrs, err := i.Addrs(); err == nil && len(addrs) > 0 {
			if i.MTU > 0 && i.MTU < mtu {
				mtu = i.MTU
			}
		}
	}
	cluster.SetMTU(mtu)
}
