package ptp

import (
	"context"
	"pan/pkg/bootstrap"
	"pan/pkg/config"
	"pan/pkg/injection"
	"pan/pkg/log"
	"pan/pkg/runtime"
	"pan/pkg/servlet"
	"reflect"
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

var (
	NetModuleName = []byte("net")
)

type stdModule struct {
	broadcastModule *stdBroadcastModule
	quicModule      *stdQuicModule

	peerGuard *stdPeerGuard

	registry runtime.Registry
	rw       sync.RWMutex
	already  bool
}

func New() interface{} {
	logger := log.Default()

	broadcastConfigurer := config.NewConfigurer[BroadcastConfig](logger)

	broadcastNetwork := &stdBroadcastNetwork{}
	broadcastNetwork.logger = logger

	broadcastServer := &stdBroadcastServer{}
	broadcastServer.logger = logger
	broadcastServer.reloadChan = make(chan struct{}, 1)

	broadcastModule := &stdBroadcastModule{}
	broadcastModule.broadcastServer = broadcastServer
	broadcastModule.broadcastNetwork = broadcastNetwork
	broadcastModule.broadcastConfigurer = broadcastConfigurer

	quicConfigurer := config.NewConfigurer[QuicConfig](logger)

	peerGuard := &stdPeerGuard{}

	quicServer := &stdQuicServer{}
	quicServer.logger = logger
	quicServer.reloadChan = make(chan struct{}, 1)
	quicServer.guard = peerGuard

	quicNetwork := &stdQuicNetwork{}
	quicNetwork.logger = logger
	quicNetwork.server = quicServer
	quicNetwork.guard = peerGuard

	quicServer.network = quicNetwork

	quicAgent := &stdQuicAgent{}
	quicAgent.logger = logger
	quicAgent.quicNetwork = quicNetwork
	quicAgent.server = quicServer
	quicAgent.broadcastNetwork = broadcastNetwork
	quicAgent.broadcastCache = expirable.NewLRU[string, uint64](0, nil, time.Second*30)
	quicAgent.guard = peerGuard
	quicAgent.reloadChan = make(chan struct{}, 1)

	quicModule := &stdQuicModule{}
	quicModule.quicServer = quicServer
	quicModule.quicNetwork = quicNetwork
	quicModule.quicAgent = quicAgent
	quicModule.quicConfigurer = quicConfigurer

	module := &stdModule{}
	module.broadcastModule = broadcastModule
	module.quicModule = quicModule
	module.peerGuard = peerGuard

	return module
}

var _ = (runtime.InitializeModule)((*stdModule)(nil))

func (module *stdModule) Init(ctx context.Context, registry runtime.Registry) error {
	module.rw.Lock()
	module.registry = registry
	module.rw.Unlock()

	var err error
	if module.already {
		err = module.reloadGuardGuides(ctx, registry)
		if err == nil {
			err = module.broadcastModule.reloadModules(ctx, registry)
		}
		if err == nil {
			err = module.quicModule.reloadModules(ctx, registry)
		}
		if err == nil {
			err = module.quicModule.reloadGuides(ctx, registry)
		}
		if err == nil {
			err = module.quicModule.reloadListeners(ctx, registry)
		}
	}
	return err
}

var _ = (bootstrap.DeferModule)((*stdModule)(nil))

func (module *stdModule) Defer(ctx context.Context) error {
	if module.already {
		return nil
	}
	module.already = true

	module.rw.RLock()
	registry := module.registry
	module.rw.RUnlock()

	var err error
	err = module.reloadGuardGuides(ctx, registry)
	if err == nil {
		err = module.broadcastModule.reloadModules(ctx, registry)
	}
	if err == nil {
		err = module.quicModule.reloadModules(ctx, registry)
	}
	if err == nil {
		err = module.quicModule.reloadGuides(ctx, registry)
	}
	if err == nil {
		err = module.quicModule.reloadListeners(ctx, registry)
	}
	return err
}

var _ = (runtime.ProviderModule)((*stdModule)(nil))

func (module *stdModule) Modules() []interface{} {
	return []interface{}{
		module.broadcastModule,
		module.quicModule,
	}
}

var _ = (runtime.EngineExtensionModule)((*stdModule)(nil))

func (module *stdModule) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[PeerGuardGuide](),
	}
}

func (module *stdModule) reloadGuardGuides(ctx context.Context, registry runtime.Registry) error {
	if ctxErr := runtime.EnsureContext(ctx); ctxErr != nil {
		return ctxErr
	}
	guides := runtime.ModulesForType[PeerGuardGuide](registry)
	module.peerGuard.setup(guides)
	return nil
}

type stdBroadcastModule struct {
	broadcastServer     *stdBroadcastServer
	broadcastNetwork    *stdBroadcastNetwork
	broadcastConfigurer BroadcastConfigurer
}

var _ = (injection.ComponentProvider)((*stdBroadcastModule)(nil))

func (module *stdBroadcastModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent[BroadcastConfigurer](module.broadcastConfigurer, injection.ComponentExternalScope),
	}
}

var _ = (BroadcastConfigListener)((*stdBroadcastModule)(nil))

func (module *stdBroadcastModule) OnConfigUpdated(config BroadcastConfig) {
	if config == nil {
		return
	}

	module.broadcastNetwork.setup(config)
	module.broadcastServer.setup(config)
}

var _ = (runtime.EngineExtensionModule)((*stdBroadcastModule)(nil))

func (module *stdBroadcastModule) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[BroadcastServeModule](),
	}
}

func (module *stdBroadcastModule) reloadModules(ctx context.Context, registry runtime.Registry) error {
	if ctxErr := runtime.EnsureContext(ctx); ctxErr != nil {
		return ctxErr
	}
	broadcastServeModules := runtime.ModulesForType[BroadcastServeModule](registry)
	module.broadcastServer.setupModules(broadcastServeModules)
	return nil
}

var _ = (bootstrap.ReadyModule)((*stdBroadcastModule)(nil))

func (module *stdBroadcastModule) Ready(ctx context.Context) error {
	module.broadcastConfigurer.Subscribe(module)
	defer module.broadcastConfigurer.Unsubscribe(module)

	return module.broadcastServer.listenAndServe(ctx)
}

type stdQuicModule struct {
	quicServer     *stdQuicServer
	quicNetwork    *stdQuicNetwork
	quicAgent      *stdQuicAgent
	quicConfigurer QuicConfigurer
}

var _ = (injection.ComponentProvider)((*stdQuicModule)(nil))

func (module *stdQuicModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent[QuicConfigurer](module.quicConfigurer, injection.ComponentExternalScope),
		injection.NewComponent[PeerNetwork](module.quicNetwork, injection.ComponentExternalScope),
		injection.NewComponent[PeerServer](module.quicServer, injection.ComponentExternalScope),
	}
}

var _ = (QuicConfigListener)((*stdQuicModule)(nil))

func (module *stdQuicModule) OnConfigUpdated(config QuicConfig) {
	if config == nil {
		return
	}

	module.quicServer.setup(config)
	module.quicNetwork.setup(config)
	module.quicAgent.setup(config)
}

var _ = (runtime.EngineExtensionModule)((*stdQuicModule)(nil))

func (module *stdQuicModule) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[PeerTopicProvider](),
		reflect.TypeFor[PeerAppModule](),
		reflect.TypeFor[QuicGuide](),
		reflect.TypeFor[PeerServerListener](),
	}
}

func (module *stdQuicModule) reloadModules(ctx context.Context, registry runtime.Registry) error {

	peerServlet := servlet.NewServlet[PeerServletContext]()

	err := runtime.TraverseRegistry(registry, func(module PeerAppModule) error {
		if ctxErr := runtime.EnsureContext(ctx); ctxErr != nil {
			return ctxErr
		}
		return module.SetupToPeer(peerServlet)
	})

	if err == nil {
		err = runtime.TraverseRegistry(registry, func(module PeerTopicProvider) error {
			if ctxErr := runtime.EnsureContext(ctx); ctxErr != nil {
				return ctxErr
			}
			topics := module.PeerTopics()
			if len(topics) <= 0 {
				return nil
			}
			var router PeerServletRouter
			router = peerServlet
			if module, ok := module.(PeerRouteModule); ok {
				routeScope := module.PeerRouteScope()
				if len(routeScope) > 0 {
					router = peerServlet.Route(routeScope)
				}
			}

			for _, topic := range topics {
				err := topic.SetupToPeer(router)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}

	if err == nil {
		module.quicServer.setupPeerServlet(peerServlet)
	}
	return err
}

func (module *stdQuicModule) reloadGuides(ctx context.Context, registry runtime.Registry) error {
	if ctxErr := runtime.EnsureContext(ctx); ctxErr != nil {
		return ctxErr
	}
	guides := runtime.ModulesForType[QuicGuide](registry)
	module.quicAgent.quicNetwork.setupGuides(guides)
	return nil
}

func (module *stdQuicModule) reloadListeners(ctx context.Context, registry runtime.Registry) error {
	if ctxErr := runtime.EnsureContext(ctx); ctxErr != nil {
		return ctxErr
	}
	listeners := runtime.ModulesForType[PeerServerListener](registry)
	module.quicServer.setupListeners(listeners)
	return nil
}

var _ = (bootstrap.ReadyModule)((*stdQuicModule)(nil))

func (module *stdQuicModule) Ready(ctx context.Context) error {

	module.quicConfigurer.Subscribe(module)
	defer module.quicConfigurer.Unsubscribe(module)

	var wg sync.WaitGroup
	wg.Add(3)

	go func(server *stdQuicServer) {
		defer wg.Done()
		server.listenAndServe(ctx)
	}(module.quicServer)

	go func(agent *stdQuicAgent) {
		defer wg.Done()
		agent.deliverBroadcast(ctx)
	}(module.quicAgent)

	go func(agent *stdQuicAgent) {
		defer wg.Done()
		agent.optimizeNetwork(ctx)
	}(module.quicAgent)

	wg.Wait()
	return nil
}

var _ = (runtime.ProviderModule)((*stdQuicModule)(nil))

func (module *stdQuicModule) Modules() []interface{} {
	return []interface{}{
		module.quicAgent,
	}
}

var _ = (PeerRouteModule)((*stdQuicModule)(nil))

func (module *stdQuicModule) PeerRouteScope() PeerServletScope {
	return NetModuleName
}
