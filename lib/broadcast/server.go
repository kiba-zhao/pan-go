// Define broadcast server for broadcast
package broadcast

import (
	"bytes"
	"context"
	"errors"
	"iter"
	"net"
	"slices"
	"strconv"
	"sync"
	"time"

	"pan/lib/log"
	libNet "pan/lib/net"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var ErrBroadcastServerUnavailable = errors.New("broadcast.BroadcastServer Error: Unavailable")
var ErrBroadcastInvalidConn = errors.New("broadcast.BroadcastServer Error: Invalid Connection")

type PacketConn interface {
	JoinGroup(ifi *net.Interface, group net.Addr) error
	LeaveGroup(ifi *net.Interface, group net.Addr) error
}

type stdBroadcastServer struct {
	logger log.Logger

	cluster *stdBroadcastCluster

	reloadChan chan struct{}
	reloadLock sync.Mutex
	reload     bool

	addrs   []string
	addrsRW sync.RWMutex
}

func (server *stdBroadcastServer) Addrs() []string {
	server.addrsRW.RLock()
	defer server.addrsRW.RUnlock()
	return server.addrs
}

func (server *stdBroadcastServer) SetAddrs(addrs []string) {
	server.logger.Debug("BroadcastServer", "SetAddrs")

	server.addrsRW.Lock()
	defer server.addrsRW.Unlock()
	if slices.Equal(server.addrs, addrs) {
		return
	}
	server.addrs = addrs
	server.Reload()
}

func (server *stdBroadcastServer) Reload() {
	server.logger.Debug("BroadcastServer", "Reload")

	server.reloadLock.Lock()
	defer server.reloadLock.Unlock()
	if server.reload {
		return
	}

	server.reload = true
	server.reloadChan <- struct{}{}
}

func (server *stdBroadcastServer) ListenAndServe(ctx context.Context) error {

	server.logger.Debug("BroadcastServer", "ListenAndServe begin")
	defer server.logger.Debug("BroadcastServer", "ListenAndServe end")

	var err error
	var closed bool

	var wg sync.WaitGroup
	var connections []*net.UDPConn

	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			closed = true
		case <-server.reloadChan:
			server.reloadLock.Lock()
			server.reload = false
			server.reloadLock.Unlock()
		}

		if len(connections) > 0 {
			for _, conn := range connections {
				conn.Close()
			}
			wg.Wait()
		}

		if closed {
			break
		}

		addrs := seqForMulitcastAddrs(server.logger, server.Addrs())
		if addrs == nil {
			continue
		}

		connections = make([]*net.UDPConn, 0)
		for port, groups := range addrs {
			conn, err := newMulticastConn(port, groups)
			if err != nil {
				server.logger.Error("BroadcastServer", "Conn Error: "+err.Error())
				continue
			}
			connections = append(connections, conn)
			wg.Add(1)
			go func(conn *net.UDPConn) {
				defer wg.Done()
				err := server.serve(conn)
				if err != nil {
					server.logger.Error("BroadcastServer", "Serve Error: "+err.Error())
				}
			}(conn)
		}
	}
	return err
}

func (server *stdBroadcastServer) serve(conn *net.UDPConn) error {
	cluster := server.cluster
	if cluster == nil {
		return ErrBroadcastServerUnavailable
	}
	mtu := cluster.MTU()
	if mtu <= 0 {
		return ErrBroadcastServerUnavailable
	}

	packetBuffers := make([]*stdPacketBuffer, 0)
	var packetBuffersRW sync.RWMutex
	var err error

	for {
		block := make([]byte, mtu)
		byteLen, addr, err := conn.ReadFromUDP(block)
		if errors.Is(err, net.ErrClosed) {
			break
		}
		if err != nil {
			continue
		}

		buffer, size := parsePacketBuffer(block[:byteLen])
		packetBuffersRW.RLock()
		idx, ok := slices.BinarySearchFunc(packetBuffers, addr.String(), comparePacketBuffer)
		var bufferItem *stdPacketBuffer
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
				bufferItem = &stdPacketBuffer{
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
				go func(item *stdPacketBuffer, ctx context.Context) {
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
		go cluster.Serve(buffer, addr.String())
	}
	return err
}

func seqForMulitcastAddrs(logger log.Logger, addrs []string) iter.Seq2[int, []net.IP] {
	if len(addrs) <= 0 {
		return nil
	}

	return func(yield func(int, []net.IP) bool) {

		mAddrs := make(map[int][]net.IP)
		for _, addr := range addrs {
			var addrPort int
			host, port, err := net.SplitHostPort(addr)
			if err == nil {
				addrPort, err = strconv.Atoi(port)
			}
			if err != nil {
				logger.Warn("BroadcastServer", "Addr Error: "+err.Error())
				continue
			}
			if _, ok := mAddrs[addrPort]; !ok {
				mAddrs[addrPort] = make([]net.IP, 0)
			}
			ip := net.ParseIP(host)
			mAddrs[addrPort] = append(mAddrs[addrPort], ip)
		}

		for port, ips := range mAddrs {
			if !yield(port, ips) {
				break
			}
		}
	}
}

func newMulticastConn(port int, groups []net.IP) (*net.UDPConn, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	addrStat, err := libNet.StatAddr()
	if err != nil {
		return nil, err
	}

	var mAddr *net.UDPAddr
	if addrStat.IPv6Enabled {
		mAddr = &net.UDPAddr{IP: net.IPv6zero, Port: port}
	} else {
		mAddr = &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: port}
	}

	conn, err := net.ListenUDP("udp", mAddr)
	if err != nil {
		return nil, err
	}

	var ipv4Conn *ipv4.PacketConn
	var ipv6Conn *ipv6.PacketConn
	for _, iface := range ifaces {

		if net.FlagMulticast != (net.FlagMulticast & iface.Flags) {
			continue
		}
		if net.FlagRunning != (net.FlagRunning & iface.Flags) {
			continue
		}

		for _, group := range groups {
			var packetConn PacketConn
			var mIP net.IP
			if mIP = group.To4(); mIP != nil {
				if ipv4Conn == nil {
					ipv4Conn = ipv4.NewPacketConn(conn)
				}
				packetConn = ipv4Conn
			} else {
				if ipv6Conn == nil {
					ipv6Conn = ipv6.NewPacketConn(conn)
				}
				packetConn = ipv6Conn
				mIP = mAddr.IP
			}

			packetConn.JoinGroup(&iface, &net.UDPAddr{IP: mIP})
		}

	}

	return conn, nil
}
