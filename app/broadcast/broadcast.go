package broadcast

import (
	"bytes"
	"context"
	"errors"
	"net"
	"pan/app/bootstrap"
	"pan/app/config"
	"pan/runtime"
	"reflect"
	"slices"
	"sync"
)

var ErrBroadcastModuleServeUnavailable = errors.New("broadcast.BroadcastModule Error: Serve Unavailable")

type BroadcastServeModule interface {
	ServeBroadcast([]byte, string) error
}

type BroadcastModule interface {
	Serve([]byte, string) error
	Deliver([]byte) error
}

func New() BroadcastModule {
	return &broadcastModule{}
}

type broadcastModule struct {
	registry       runtime.Registry
	registryLocker sync.RWMutex

	addresses  []string
	locker     sync.RWMutex
	reloadChan chan struct{}
	reloadOnce sync.Once
	needReload bool
	mtu        int
}

func (b *broadcastModule) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[BroadcastServeModule](),
	}
}

func (b *broadcastModule) Components() []bootstrap.Component {
	return []bootstrap.Component{
		bootstrap.NewComponent[BroadcastModule](b, bootstrap.ComponentExternalScope),
	}
}

func (b *broadcastModule) Addresses() []string {
	b.locker.RLock()
	defer b.locker.RUnlock()
	return b.addresses
}

func (b *broadcastModule) Serve(payload []byte, ip string) error {
	b.registryLocker.RLock()
	registry := b.registry
	b.registryLocker.RUnlock()

	if registry == nil {
		return ErrBroadcastModuleServeUnavailable
	}

	return runtime.TraverseRegistry(registry, func(module BroadcastServeModule) error {
		return module.ServeBroadcast(payload, ip)
	})
}

func (b *broadcastModule) Deliver(payload []byte) error {
	size := len(payload)
	if size <= 0 {
		return nil
	}
	if size > 65531 {
		return bytes.ErrTooLarge
	}

	addresses := b.Addresses()

	connArr := make([]*net.UDPConn, 0)
	for _, address := range addresses {
		addr, err := net.ResolveUDPAddr("udp", address)
		if err != nil {
			return err
		}
		conn, err := net.DialUDP("udp", nil, addr)
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

func (b *broadcastModule) ReloadChan() chan struct{} {

	b.reloadOnce.Do(func() {
		b.reloadChan = make(chan struct{}, 1)
		b.mtu = broadcastMTU()
	})

	return b.reloadChan
}

func (b *broadcastModule) OnConfigUpdated(settings config.AppSettings) {
	b.locker.Lock()
	defer b.locker.Unlock()

	if slices.Equal(b.addresses, settings.BroadcastAddress) {
		return
	}

	b.addresses = settings.BroadcastAddress

	// trigger to reload
	if b.needReload {
		return
	}
	b.needReload = true
	b.ReloadChan() <- struct{}{}
}

func (b *broadcastModule) Init(registry runtime.Registry) error {
	b.registryLocker.Lock()
	defer b.registryLocker.Unlock()
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
		addresses := b.addresses
		b.locker.Unlock()

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
		for _, address := range addresses {
			server := &broadcastServer{
				address:         address,
				broadcastModule: b,
				mtu:             mtu,
			}

			servers = append(servers, server)
			wg.Add(1)
			go func(bs *broadcastServer) {
				defer wg.Done()
				_ = bs.ListenAndServe()
				// TODO: write error into log
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
