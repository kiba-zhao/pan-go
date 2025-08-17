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

	runtime *stdBroadcastRuntime
	cluster *stdBroadcastCluster
	agent   *stdBroadcastAgent
	server  *stdBroadcastServer
}

func New() interface{} {
	module := &stdBroadcastModule{}

	cluster := &stdBroadcastCluster{}
	module.cluster = cluster

	agent := &stdBroadcastAgent{}
	module.agent = agent

	agent.reloadChan = make(chan struct{}, 1)
	agent.cluster = cluster

	server := &stdBroadcastServer{}
	module.server = server
	server.reloadChan = make(chan struct{}, 1)
	server.cluster = cluster

	// init broadcast runtime
	runtime := &stdBroadcastRuntime{}
	module.runtime = runtime

	server.runtime = runtime
	runtime.server = server
	setDefualtMTU(runtime)

	agent.runtime = runtime
	runtime.agent = agent

	cluster.runtime = runtime
	runtime.RegisterServeModule(agent)
	runtime.setDeliverLimitSize(65535)
	//

	logger := log.Default()
	cluster.logger = logger
	agent.logger = logger
	server.logger = logger
	runtime.logger = logger

	return module
}

var _ = (BroadcastModule)((*stdBroadcastModule)(nil))

func (module *stdBroadcastModule) SetStore(store BroadcastStore) {
	module.runtime.setStore(store)
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
	module.runtime.setQuicCluster(module.QuicCluster)

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
	module.runtime.setAddrs(settings.BroadcastAddrs)
}

var _ = (peer.PeerConfigListener)((*stdBroadcastModule)(nil))

func (module *stdBroadcastModule) OnPeerConfigUpdated(settings *peer.PeerSettings) {
	module.runtime.setPeerSettings(settings)
}

func setDefualtMTU(runtime *stdBroadcastRuntime) {
	ifs, err := net.Interfaces()

	if err != nil {
		runtime.setMTU(1500)
		return
	}

	mtu := runtime.DeliverLimitSize()
	for _, i := range ifs {
		if addrs, err := i.Addrs(); err == nil && len(addrs) > 0 {
			if i.MTU > 0 && i.MTU < mtu {
				mtu = i.MTU
			}
		}
	}
	runtime.setMTU(mtu)
}
