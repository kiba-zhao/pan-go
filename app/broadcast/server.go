package broadcast

import (
	"bytes"
	"context"
	"errors"
	"net"
	"slices"
	"sync"
	"time"
)

var ErrBroadcastServerUnavailable = errors.New("broadcast.BroadcastServer Error: Unavailable")
var ErrBroadcastInvalidConn = errors.New("broadcast.BroadcastServer Error: Invalid Connection")

type broadcastServer struct {
	broadcastModule *broadcastModule
	locker          sync.RWMutex
	conn            *net.UDPConn
	address         string
	mtu             int
}

func (bs *broadcastServer) Shutdown() error {
	bs.locker.RLock()
	conn := bs.conn
	bs.locker.RUnlock()
	if conn == nil {
		return ErrBroadcastInvalidConn
	}

	bs.locker.Lock()
	bs.conn = nil
	bs.locker.Unlock()
	return conn.Close()
}

func (bs *broadcastServer) ListenAndServe() error {
	if bs.broadcastModule == nil {
		return ErrBroadcastServerUnavailable
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

	// bucket_ := cache.NewBucket[string, *PacketBuffer](cmp.Compare[string])
	// bufferBucket := cache.WrapSyncBucket(bucket_)
	packetBuffers := make([]*PacketBuffer, 0)
	var packetBuffersRW sync.RWMutex

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
		packetBuffersRW.RLock()
		// bufferItem, ok := bufferBucket.Search(addr.String())
		idx, ok := slices.BinarySearchFunc(packetBuffers, addr.String(), comparePacketBuffer)
		var bufferItem *PacketBuffer
		if ok {
			bufferItem = packetBuffers[idx]
		}
		packetBuffersRW.RUnlock()

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
				bufferItem = &PacketBuffer{
					addr:   addr.String(),
					cancel: cancel,
				}
				bufferItem.wg.Add(1)
				packetBuffersRW.Lock()
				idx_, ok_ := slices.BinarySearchFunc(packetBuffers, addr.String(), comparePacketBuffer)
				if !ok_ {
					packetBuffers = slices.Insert(packetBuffers, idx_, bufferItem)
				}
				packetBuffersRW.Unlock()
				go func(item *PacketBuffer, ctx context.Context) {
					defer item.wg.Done()
					<-ctx.Done()
					packetBuffersRW.Lock()
					defer packetBuffersRW.Unlock()
					idx_, ok_ := slices.BinarySearchFunc(packetBuffers, item.addr, comparePacketBuffer)
					if ok_ {
						packetBuffers = slices.Delete(packetBuffers, idx_, idx_+1)
					}
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
		go bs.broadcastModule.Serve(buffer, addr.IP.String())
	}
	return err
}

func (b *broadcastServer) HashCode() string {
	return b.address
}
