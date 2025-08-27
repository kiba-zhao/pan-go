package broadcast

import (
	"context"
	"pan/lib/bootstrap"
	"pan/lib/config"
	"pan/lib/injection"
	"pan/lib/log"
	"pan/lib/quic"
	"sync"
)

type BroadcastModule interface {
	SetStore(store BroadcastStore)
}

type stdBroadcastModule struct {
	QuicNetwork quic.QuicNetwork

	configurer BroadcastConfigurer
	network    *stdBroadcastNetwork
	agent      *stdBroadcastAgent
	server     *stdBroadcastServer
	provider   *stdBroadcastProvider
}

func New() interface{} {

	network := &stdBroadcastNetwork{}

	agent := &stdBroadcastAgent{}
	agent.network = network
	agent.reloadChan = make(chan struct{}, 1)

	server := &stdBroadcastServer{}
	server.network = network
	server.reloadChan = make(chan struct{}, 1)

	provider := &stdBroadcastProvider{}
	network.provider = provider
	agent.provider = provider

	logger := log.Default()
	network.logger = logger
	agent.logger = logger
	server.logger = logger

	configurer := config.NewConfigurer[BroadcastConfig](logger)

	module := &stdBroadcastModule{}
	module.network = network
	module.agent = agent
	module.server = server
	module.provider = provider
	module.configurer = configurer

	return module
}

var _ = (BroadcastModule)((*stdBroadcastModule)(nil))

func (module *stdBroadcastModule) SetStore(store BroadcastStore) {
	module.agent.SetupStore(store)
}

var _ = (injection.ComponentProvider)((*stdBroadcastModule)(nil))

func (module *stdBroadcastModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(module, injection.ComponentNoneScope),
		injection.NewComponent[BroadcastModule](module, injection.ComponentExternalScope),
		injection.NewComponent[BroadcastNetwork](module.network, injection.ComponentExternalScope),
		injection.NewComponent[BroadcastConfigurer](module.configurer, injection.ComponentExternalScope),
	}
}

var _ = (bootstrap.ReadyModule)((*stdBroadcastModule)(nil))

func (module *stdBroadcastModule) Ready(ctx context.Context) error {
	module.provider.SetupQuicNetwork(module.QuicNetwork)

	module.configurer.Subscribe(module)
	defer module.configurer.Unsubscribe(module)

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

var _ = (BroadcastConfigListener)((*stdBroadcastModule)(nil))

func (module *stdBroadcastModule) OnConfigUpdated(config BroadcastConfig) {
	module.network.Setup(config)
	module.server.Setup(config)
	module.agent.Setup(config)
}
