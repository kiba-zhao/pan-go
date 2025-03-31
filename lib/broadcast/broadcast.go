// Package broadcast provides the broadcast engine
package broadcast

import (
	"bytes"
	"context"
	"errors"
	"net"
	"pan/lib/injection"
	"pan/lib/runtime"
	"reflect"
	"slices"
	"sync"
)

var ErrBroadcastModuleServeUnavailable = errors.New("broadcast.BroadcastModule Error: Serve Unavailable")
var ErrBroadcastDeliverNoAddrs = errors.New("broadcast.BroadcastModule Error: Deliver No Addrs")

// BroadcastServeAddrsProvider is the broadcast serve addrs provider
type BroadcastServeAddrsProvider interface {
	// BroadcastServeAddrs returns the list of addresses that the broadcast module can serve.
	BroadcastServeAddrs() []string
}

// BroadcastDeliverAddrsProvider is the broadcast deliver addrs provider
type BroadcastDeliverAddrsProvider interface {
	// BroadcastDeliverAddrs returns the list of addresses that the broadcast module can deliver.
	BroadcastDeliverAddrs() []string
}

// BroadcastPublicAddrsProvider is the broadcast public addrs provider
type BroadcastPublicAddrsProvider interface {
	// BroadcastPublicAddrs returns the list of addresses that the broadcast module can deliver.
	BroadcastPublicAddrs() []string
}

// BroadcastServeModule is the broadcast serve module
type BroadcastServeModule interface {
	// ServeBroadcast serves the broadcast message to the peer.
	ServeBroadcast([]byte, string) error
}

// BroadcastModule is the broadcast module
type BroadcastModule interface {
	// Serve handles the broadcasting of a message.
	// Returns an error if the broadcast module is unavailable or if there is an issue serving the broadcast.
	Serve([]byte, string) error
	// Deliver broadcasts a message to the peers.
	//
	// The message is delivered to the peers whose addresses are in the parameter list.
	// If the parameter list is empty, the message is delivered to all online peers.
	//
	// Returns an error if the broadcast module is unavailable or if there is an issue delivering the broadcast.
	Deliver([]byte, ...string) error
	// ServeAddrs returns the list of addresses that the broadcast module can serve.
	//
	// This function is called by the runtime to get the list of addresses that the
	// broadcast module can serve. The addresses are used to serve the broadcast
	// message when the Serve method is called.
	ServeAddrs() []string
	// DeliverAddrs returns the list of addresses that the broadcast module can deliver.
	//
	// The addresses are used to deliver the broadcast message when the Deliver method is called.
	// If the parameter list is empty, the message is delivered to all online peers.
	DeliverAddrs() []string
	// Reload reloads the broadcast module.
	//
	// It is called by the runtime to reload the broadcast module after the application has finished initializing.
	Reload()
}

// New creates a new broadcast module.
//
// The module implements the BroadcastModule interface and is used to manage broadcast messages.
//
// The module is also a runtime.Module, and can be used to load components into the broadcast engine.
//
// The module is initialized with the given ComponentStoreProvider, which is used to get the ComponentStore.
// The ComponentStore is used to store components that are injected into other components.
//
// The module is responsible for managing the broadcast engine and providing the necessary methods to serve and deliver broadcast messages.
func New(csProvider injection.ComponentStoreProvider) BroadcastModule {
	module := &broadcastModule{}

	agent := &broadcastAgent{provider: csProvider}
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
	b.ReloadChan() <- struct{}{}
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
		udpAddrs, err := resolveAddrs(addr)
		if err != nil {
			return err
		}
		if len(udpAddrs) <= 0 {
			continue
		}
		for _, udpAddr := range udpAddrs {
			conn, err := net.DialUDP("udp", nil, udpAddr)
			if err != nil {
				continue
			}
			connArr = append(connArr, conn)
			defer conn.Close()
		}
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

	var server broadcastServer
	server.module = b
	server.mtu = b.mtu

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

		if closed {
			server.Shutdown()
			break
		}

		server.ListenAndServe(addrs)
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
		if addrs, err := i.Addrs(); err == nil && len(addrs) > 0 {
			if i.MTU > 0 && i.MTU < mtu {
				mtu = i.MTU
			}
		}
	}
	return mtu
}

func resolveAddrs(addr string) ([]*net.UDPAddr, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	if udpAddr.IP.To4() != nil || !udpAddr.IP.IsMulticast() || len(udpAddr.Zone) > 0 {
		return []*net.UDPAddr{udpAddr}, nil
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	isGlobalAddr := isGlobalMulticastIP(udpAddr.IP)
	udpAddrs := make([]*net.UDPAddr, 0)
	for _, iface := range ifaces {
		if net.FlagMulticast != (net.FlagMulticast & iface.Flags) {
			continue
		}
		if net.FlagRunning != (net.FlagRunning & iface.Flags) {
			continue
		}
		if net.FlagLoopback == (net.FlagLoopback & iface.Flags) {
			continue
		}
		ifaceAddrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, ifaceAddr := range ifaceAddrs {
			ipNet, ok := ifaceAddr.(*net.IPNet)
			if !ok || ipNet.IP.To4() != nil {
				continue
			}
			if isGlobalAddr == (ipNet.IP.IsGlobalUnicast() && !ipNet.IP.IsPrivate()) {
				udpAddr_ := net.UDPAddrFromAddrPort(udpAddr.AddrPort())
				udpAddr_.Zone = iface.Name
				udpAddrs = append(udpAddrs, udpAddr_)
				break
			}
		}
	}
	return udpAddrs, err
}

func isGlobalMulticastIP(ip net.IP) bool {
	return !(ip.IsLinkLocalMulticast() || ip.IsInterfaceLocalMulticast())
}
