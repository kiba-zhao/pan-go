package app

import (
	"bytes"
	"cmp"
	"context"
	"encoding/binary"
	"errors"
	"net"
	"pan/app/cache"
	"pan/runtime"
	"reflect"
	"slices"
	"sync"
	"time"
)

type BroadcastModule interface {
	ServeBroadcast([]byte) error
}

type broadcastPacketBuffer struct {
	size    int
	content []byte
	addr    string
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func (bpb *broadcastPacketBuffer) HashCode() string {
	return bpb.addr
}

type broadcastServer struct {
	broadcast *broadcast
	locker    sync.RWMutex
	conn      *net.UDPConn
	address   string
	mtu       int
}

func (bs *broadcastServer) Shutdown() error {
	bs.locker.RLock()
	conn := bs.conn
	bs.locker.RUnlock()
	if conn == nil {
		return ErrUnavailable
	}

	bs.locker.Lock()
	bs.conn = nil
	bs.locker.Unlock()
	return conn.Close()
}

func (bs *broadcastServer) ListenAndServe() error {
	if bs.broadcast == nil {
		return ErrUnavailable
	}

	addr, err := net.ResolveUDPAddr("udp", bs.address)
	if err != nil {
		return err
	}

	conn, err := net.ListenMulticastUDP(addr.Network(), nil, addr)
	if err != nil {
		return err
	}
	bs.locker.Lock()
	bs.conn = conn
	bs.locker.Unlock()
	defer bs.Shutdown()

	bucket_ := cache.NewBucket[string, *broadcastPacketBuffer](cmp.Compare[string])
	bufferBucket := cache.WrapSyncBucket(bucket_)

	for {
		block := make([]byte, bs.mtu)
		byteLen, addr, err := conn.ReadFromUDP(block)
		if errors.Is(err, net.ErrClosed) {
			break
		}
		if err != nil {
			continue
		}

		buffer, size := parsePacketBuffer(block[:byteLen])
		bufferItem, ok := bufferBucket.Search(addr.String())
		if size > 0 && ok {
			bufferItem.cancel()
			bufferItem.wg.Wait()
			ok = false
		} else if size == 0 {
			if !ok {
				continue
			}
			len_ := len(bufferItem.content) + len(block)
			if len_ > bufferItem.size {
				continue
			}
			buffer = bytes.Join([][]byte{bufferItem.content, block}, nil)
			size = bufferItem.size
		}

		if size > len(buffer) {
			if !ok {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				bufferItem = &broadcastPacketBuffer{
					addr:   addr.String(),
					cancel: cancel,
				}
				bufferItem.wg.Add(1)
				bufferBucket.Store(bufferItem)
				go func(item *broadcastPacketBuffer, ctx context.Context) {
					defer item.wg.Done()
					<-ctx.Done()
					bufferBucket.Delete(item)
				}(bufferItem, ctx)
			}

			bufferItem.size = size
			bufferItem.content = buffer
			continue
		}

		if ok {
			bufferItem.cancel()
			bufferItem.wg.Wait()
		}
		go bs.broadcast.Serve(buffer)
	}
	return err
}

func (b *broadcastServer) HashCode() string {
	return b.address
}

func parsePacketBuffer(block []byte) ([]byte, int) {
	if block[0] != 0 {
		return block, 0
	}
	offset := 1
	checksum := binary.BigEndian.Uint16(block[offset : offset+2])
	offset += 2
	size16 := binary.BigEndian.Uint16(block[offset : offset+2])
	offset += 2
	if checksum^size16 != 0 {
		return block, 0
	}

	size := int(size16)
	if size < len(block)-4 {
		return block, 0
	}
	return block[offset:], size
}

func packBuffer(buffer []byte) []byte {
	size := len(buffer)
	checksum := 65535 ^ uint16(size)

	return bytes.Join([][]byte{
		binary.BigEndian.AppendUint16([]byte{0}, checksum),
		binary.BigEndian.AppendUint16(nil, uint16(size)),
		buffer,
	}, nil)
}

type Broadcast interface {
	Serve([]byte) error
	Deliver([]byte) error
}

type broadcast struct {
	registry       runtime.Registry
	registryLocker sync.RWMutex

	addresses []string
	locker    sync.RWMutex
	sigChan   chan bool
	sigOnce   sync.Once
	hasSig    bool
	mtu       int
}

func (b *broadcast) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[BroadcastModule](),
	}
}

func (b *broadcast) Components() []runtime.Component {
	return []runtime.Component{
		runtime.NewComponent[Broadcast](b, runtime.ComponentExternalScope),
	}
}

func (b *broadcast) Addresses() []string {
	b.locker.RLock()
	defer b.locker.RUnlock()
	return b.addresses
}

func (b *broadcast) Serve(payload []byte) error {
	b.registryLocker.RLock()
	registry := b.registry
	b.registryLocker.RUnlock()

	if registry == nil {
		return ErrUnavailable
	}

	return runtime.TraverseRegistry(registry, func(module BroadcastModule) error {
		return module.ServeBroadcast(payload)
	})
}

func (b *broadcast) Deliver(payload []byte) error {
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

func (b *broadcast) setSig(sig bool) {
	if b.hasSig {
		return
	}

	b.sigOnce.Do(func() {
		b.sigChan = make(chan bool, 1)
		b.mtu = broadcastMTU()
	})

	b.hasSig = true
	b.sigChan <- sig
}

func (b *broadcast) OnConfigUpdated(settings AppSettings) {
	b.locker.Lock()
	defer b.locker.Unlock()

	if slices.Equal(b.addresses, settings.BroadcastAddress) {
		return
	}

	b.addresses = settings.BroadcastAddress
	b.setSig(true)
}

func (b *broadcast) Init(registry runtime.Registry) error {
	b.registryLocker.Lock()
	defer b.registryLocker.Unlock()
	b.registry = registry
	return nil
}

func (b *broadcast) Ready() error {

	var serverMgr cache.Bucket[string, *broadcastServer]
	b.sigOnce.Do(func() {
		b.sigChan = make(chan bool, 1)
		b.mtu = broadcastMTU()
	})

	for {
		sig := <-b.sigChan
		b.locker.Lock()
		b.hasSig = false
		addresses := b.addresses
		b.locker.Unlock()

		if serverMgr != nil {
			for _, item := range serverMgr.Items() {
				item.Shutdown()
			}
		}

		if !sig {
			break
		}

		serverMgr = cache.NewBucket[string, *broadcastServer](cmp.Compare[string])
		mtu := b.mtu
		for _, address := range addresses {
			server := &broadcastServer{
				address:   address,
				broadcast: b,
				mtu:       mtu,
			}

			serverMgr.Store(server)
			go func(bs *broadcastServer) {
				for {
					err := bs.ListenAndServe()
					if errors.Is(err, net.ErrClosed) {
						break
					}
					time.Sleep(6 * time.Second)
				}
			}(server)
		}
	}

	return nil
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
