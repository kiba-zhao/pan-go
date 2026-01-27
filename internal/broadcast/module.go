package broadcast

import (
	"context"
	"pan/internal/bootstrap"
	"pan/internal/config"
	"pan/internal/injection"
	"pan/internal/log"
	"pan/internal/quic"
	"pan/internal/repository"
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

type stdBroadcastModule struct {
	QuicNetwork       quic.QuicNetwork
	RepositoryManager repository.RepositoryManager

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
	agent.cache = expirable.NewLRU[string, uint64](2000, nil, time.Second*30)

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

var _ = (injection.ComponentProvider)((*stdBroadcastModule)(nil))

func (module *stdBroadcastModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(module, injection.ComponentNoneScope),
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
