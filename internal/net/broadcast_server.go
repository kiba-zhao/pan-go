package net

import (
	"bytes"
	"context"
	"errors"
	"iter"
	"net"
	"pan/internal/log"
	"slices"
	"strconv"
	"sync"
	"time"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var ErrBroadcastServerUnavailable = errors.New("net.BroadcastServer Error: Unavailable")
var ErrBroadcastInvalidConn = errors.New("net.BroadcastServer Error: Invalid Connection")

type PacketConn interface {
	JoinGroup(ifi *net.Interface, group net.Addr) error
	LeaveGroup(ifi *net.Interface, group net.Addr) error
}

type BroadcastServeModule interface {
	// ServeBroadcast serves the broadcast message to the peer.
	ServeBroadcast([]byte, string) error
}

type stdBroadcastServer struct {
	logger log.Logger

	addrs       []string
	mtu         int
	ipv6Enabled bool

	reloadChan chan struct{}
	reloadLock sync.Mutex
	reload     bool

	modules  []BroadcastServeModule
	moduleRW sync.RWMutex
}

func (server *stdBroadcastServer) SetupModules(modules []BroadcastServeModule) {
	server.moduleRW.Lock()
	defer server.moduleRW.Unlock()
	server.modules = modules
}

func (server *stdBroadcastServer) Setup(config BroadcastConfig) {
	server.logger.Debug("net.BroadcastServer", "Setup")

	server.reloadLock.Lock()
	defer server.reloadLock.Unlock()

	changed := false
	addrs := config.Addrs()
	mtu := config.MTU()
	ipv6Enabled := len(config.IPv6ZoneList()) > 0

	if !slices.Equal(server.addrs, addrs) {
		server.addrs = addrs
		changed = true
	}

	if server.mtu != mtu {
		server.mtu = mtu
		changed = true
	}

	if server.ipv6Enabled != ipv6Enabled {
		server.ipv6Enabled = ipv6Enabled
		changed = true
	}

	if !changed || server.reload {
		return
	}

	server.reload = true
	server.reloadChan <- struct{}{}
}

func (server *stdBroadcastServer) ListenAndServe(ctx context.Context) error {

	server.logger.Debug("net.BroadcastServer", "ListenAndServe begin")
	defer server.logger.Debug("net.BroadcastServer", "ListenAndServe end")

	var err error
	var wg sync.WaitGroup
	var connections []*net.UDPConn
	var addrs []string
	var mtu int
	var ipv6Enabled bool

	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		case <-server.reloadChan:
			server.reloadLock.Lock()
			server.reload = false
			addrs = server.addrs
			mtu = server.mtu
			ipv6Enabled = server.ipv6Enabled
			server.reloadLock.Unlock()
		}

		if len(connections) > 0 {
			for _, conn := range connections {
				conn.Close()
			}
			wg.Wait()
		}

		if err != nil {
			break
		}

		if len(addrs) <= 0 || mtu < 0 {
			continue
		}

		mAddrs := seqForMulitcastAddrs(server.logger, addrs)
		if addrs == nil {
			continue
		}

		connections = make([]*net.UDPConn, 0)
		for port, groups := range mAddrs {
			conn, err := newMulticastConn(port, groups, server.logger, ipv6Enabled)
			if err != nil {
				server.logger.Error("net.BroadcastServer", "Conn Error: "+err.Error())
				continue
			}
			connections = append(connections, conn)
			wg.Add(1)
			go func(conn *net.UDPConn, mtu int) {
				defer wg.Done()
				err := server.Serve(conn, mtu)
				if err != nil {
					server.logger.Error("net.BroadcastServer", "Serve Error: "+err.Error())
				}
			}(conn, mtu)
		}
	}
	return err
}

func (server *stdBroadcastServer) Serve(conn *net.UDPConn, mtu int) error {

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
		go server.Handle(buffer, addr.String())
	}
	return err
}

func (server *stdBroadcastServer) Handle(payload []byte, addr string) error {
	server.logger.Debug("net.BroadcastServer", "Handle begin: "+addr)
	defer server.logger.Debug("net.BroadcastServer", "Handle end: "+addr)

	server.moduleRW.RLock()
	modules := slices.Clone(server.modules)
	server.moduleRW.RUnlock()

	var err error
	for _, module := range modules {
		err = module.ServeBroadcast(payload, addr)
		if err == nil {
			break
		}
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
				logger.Warn("net.BroadcastServer", "Addr Error: "+err.Error())
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

func newMulticastConn(port int, groups []net.IP, logger log.Logger, ipv6Enabled bool) (*net.UDPConn, error) {

	var mAddr *net.UDPAddr
	if ipv6Enabled {
		mAddr = &net.UDPAddr{IP: net.IPv6zero, Port: port}
	} else {
		mAddr = &net.UDPAddr{IP: net.IPv4zero, Port: port}
	}

	conn, err := net.ListenUDP("udp", mAddr)
	if err != nil {
		return nil, err
	}

	var ipv4Conn *ipv4.PacketConn
	var ipv6Conn *ipv6.PacketConn
	for _, group := range groups {
		var packetConn PacketConn
		var mIP net.IP
		if mIP = group.To4(); mIP != nil {
			if ipv4Conn == nil {
				ipv4Conn = ipv4.NewPacketConn(conn)
			}
			packetConn = ipv4Conn
		} else if mIP = group.To16(); mIP != nil {
			if ipv6Conn == nil {
				ipv6Conn = ipv6.NewPacketConn(conn)
			}
			packetConn = ipv6Conn
		}

		err = packetConn.JoinGroup(nil, &net.UDPAddr{IP: mIP})
		if err != nil {
			logger.Warn("net.BroadcastServer", "JoinGroup Error: "+err.Error()+" IP="+mIP.String())
		}
	}

	return conn, nil
}
