package broadcast

import (
	"bytes"
	"context"
	"errors"
	"net"
	"pan/app/injection"
	"pan/logger"
	"pan/runtime"
	"reflect"
	"slices"
	"sync"
)

var ErrBroadcastModuleServeUnavailable = errors.New("broadcast.BroadcastModule Error: Serve Unavailable")
var ErrBroadcastDeliverNoAddrs = errors.New("broadcast.BroadcastModule Error: Deliver No Addrs")

type BroadcastServeAddrsProvider interface {
	BroadcastServeAddrs() []string
}

type BroadcastDeliverAddrsProvider interface {
	BroadcastDeliverAddrs() []string
}

type BroadcastPublicAddrsProvider interface {
	BroadcastPublicAddrs() []string
}

type BroadcastServeModule interface {
	ServeBroadcast([]byte, string) error
}

type BroadcastModule interface {
	Serve([]byte, string) error
	Deliver([]byte, ...string) error
	ServeAddrs() []string
	DeliverAddrs() []string
	Reload()
}

func New(store injection.ComponentStore) BroadcastModule {
	module := &broadcastModule{}

	agent := &broadcastAgent{store: store}
	agent.module = module
	module.agent = agent

	provider := &broadcastAddrsProvider{}
	provider.agent = agent
	provider.module = module
	module.provider = provider

	return module
}

type broadcastModule struct {
	provider   *broadcastAddrsProvider
	agent      *broadcastAgent
	registry   runtime.Registry
	registryRW sync.RWMutex

	locker     sync.Mutex
	reloadChan chan struct{}
	reloadOnce sync.Once
	needReload bool
	mtu        int
}

func (b *broadcastModule) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[BroadcastServeModule](),
		reflect.TypeFor[BroadcastServeAddrsProvider](),
		reflect.TypeFor[BroadcastDeliverAddrsProvider](),
	}
}

func (b *broadcastModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent[BroadcastModule](b, injection.ComponentExternalScope),
	}
}

func (b *broadcastModule) Modules() []interface{} {
	return []interface{}{b.agent, b.provider}
}

func (b *broadcastModule) Reload() {
	b.locker.Lock()
	defer b.locker.Unlock()
	if b.needReload {
		return
	}
	b.needReload = true
	b.ReloadChan()
}

func (b *broadcastModule) Serve(payload []byte, addr string) error {

	b.registryRW.RLock()
	registry := b.registry
	b.registryRW.RUnlock()

	if registry == nil {
		return ErrBroadcastModuleServeUnavailable
	}

	return runtime.TraverseRegistry(registry, func(module BroadcastServeModule) error {
		return module.ServeBroadcast(payload, addr)
	})
}

func (b *broadcastModule) Deliver(payload []byte, addrs ...string) error {
	size := len(payload) + 1
	if size > 65531 {
		return bytes.ErrTooLarge
	}

	var deliverAddrs []string
	if len(addrs) <= 0 {
		deliverAddrs = b.DeliverAddrs()
	} else {
		deliverAddrs = addrs
	}

	if len(deliverAddrs) <= 0 {
		return ErrBroadcastDeliverNoAddrs
	}

	connArr := make([]*net.UDPConn, 0)
	for _, addr := range deliverAddrs {
		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			return err
		}
		conn, err := net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			return err
		}
		connArr = append(connArr, conn)
		defer conn.Close()
	}

	mtu := b.mtu
	buffer := packBuffer(payload)
	size = len(buffer)
	for offset := 0; offset < size; offset += mtu {
		var limit int
		if offset+mtu > size {
			limit = size - offset
		} else {
			limit = offset + mtu
		}
		block := buffer[offset:limit]

		for _, conn := range connArr {
			_, err := conn.Write(block)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *broadcastModule) ServeAddrs() []string {
	b.registryRW.RLock()
	registry := b.registry
	b.registryRW.RUnlock()

	if registry == nil {
		return nil
	}

	addrs := make([]string, 0)
	runtime.TraverseRegistry(registry, func(provider BroadcastServeAddrsProvider) error {
		serveAddrs := provider.BroadcastServeAddrs()
		if len(serveAddrs) <= 0 {
			return nil
		}
		for _, addr := range serveAddrs {
			if idx, ok := slices.BinarySearch(addrs, addr); !ok {
				addrs = slices.Insert(addrs, idx, addr)
			}
		}
		return nil
	})
	return addrs
}

func (b *broadcastModule) DeliverAddrs() []string {
	b.registryRW.RLock()
	registry := b.registry
	b.registryRW.RUnlock()

	if registry == nil {
		return nil
	}

	addrs := make([]string, 0)
	runtime.TraverseRegistry(registry, func(provider BroadcastDeliverAddrsProvider) error {
		deliverAddrs := provider.BroadcastDeliverAddrs()
		if len(deliverAddrs) <= 0 {
			return nil
		}
		for _, addr := range deliverAddrs {
			if idx, ok := slices.BinarySearch(addrs, addr); !ok {
				addrs = slices.Insert(addrs, idx, addr)
			}
		}
		return nil
	})
	return addrs
}

func (b *broadcastModule) ReloadChan() chan struct{} {

	b.reloadOnce.Do(func() {
		b.reloadChan = make(chan struct{}, 1)
		b.mtu = broadcastMTU()
	})

	return b.reloadChan
}

func (b *broadcastModule) Init(registry runtime.Registry) error {
	b.registryRW.Lock()
	defer b.registryRW.Unlock()
	b.registry = registry
	return nil
}

func (b *broadcastModule) Ready(ctx context.Context) error {

	var wg sync.WaitGroup
	var servers []*broadcastServer
	var err error
	closed := false
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			closed = true
		case <-b.ReloadChan():
		}

		b.locker.Lock()
		b.needReload = false
		b.locker.Unlock()
		addrs := b.ServeAddrs()

		if len(servers) > 0 {
			for _, item := range servers {
				item.Shutdown()
			}
			wg.Wait()
		}

		if closed {
			break
		}

		servers = make([]*broadcastServer, 0)
		mtu := b.mtu
		for _, addr := range addrs {
			server := &broadcastServer{
				address: addr,
				module:  b,
				mtu:     mtu,
			}

			servers = append(servers, server)
			wg.Add(1)
			go func(bs *broadcastServer) {
				defer wg.Done()
				err = bs.ListenAndServe()
				if err != nil {
					logger.Default().Log(context.Background(), logger.LevelError, "app.broadcast Error: %s", err.Error())
				}
			}(server)
		}
	}

	return err
}

func broadcastMTU() int {
	ifs, err := net.Interfaces()

	if err != nil {
		return 1500
	}
	mtu := 65535
	for _, i := range ifs {
		if i.MTU < mtu {
			mtu = i.MTU
		}
	}
	return mtu
}
