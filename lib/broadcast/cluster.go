package broadcast

import (
	"bytes"
	"errors"
	"net"
	"pan/lib/log"
	"slices"
	"sync"
)

var ErrBroadcastClusterServeUnavailable = errors.New("broadcast.BroadcastCluster Error: Serve Unavailable")
var ErrBroadcastClusterDeliverNoAddrs = errors.New("broadcast.BroadcastCluster Error: Deliver No Addrs")

// BroadcastServeModule is the broadcast serve module
type BroadcastServeModule interface {
	// ServeBroadcast serves the broadcast message to the peer.
	ServeBroadcast([]byte, string) error
}

type BroadcastCluster interface {
	// ServeModules returns the list of broadcast serve modules.
	ServeModules() []BroadcastServeModule
	// RegisterServeModule registers a broadcast serve module.
	RegisterServeModule(BroadcastServeModule)
	// UnregisterServeModule unregisters a broadcast serve module.
	UnregisterServeModule(BroadcastServeModule)
	// Serve serves the broadcast message to the peer.
	Serve([]byte, string) error

	// DeliverAddrs returns the list of broadcast deliver addresses.
	DeliverAddrs() []string
	// SetDeliverAddr sets the broadcast deliver addresses.
	SetDeliverAddr([]string)
	// DeliverLimitSize returns the broadcast deliver limit size.
	DeliverLimitSize() int
	// SetDeliverLimitSize sets the broadcast deliver limit size.
	SetDeliverLimitSize(int)
	// MTU returns the broadcast MTU.
	MTU() int
	// SetMTU sets the broadcast MTU.
	SetMTU(int)
	// Deliver delivers the broadcast message to the peer.
	Deliver([]byte, ...string) error
}

type stdBroadcastCluster struct {
	logger log.Logger

	serveModules   []BroadcastServeModule
	serveModulesRW sync.RWMutex

	deliverAddrs   []string
	deliverAddrsRW sync.RWMutex

	deliverLimitSize   int
	deliverLimitSizeRW sync.RWMutex

	mtu   int
	mtuRW sync.RWMutex
}

var _ = (BroadcastCluster)((*stdBroadcastCluster)(nil))

func (cluster *stdBroadcastCluster) ServeModules() []BroadcastServeModule {
	cluster.serveModulesRW.RLock()
	defer cluster.serveModulesRW.RUnlock()
	return cluster.serveModules
}

func (cluster *stdBroadcastCluster) RegisterServeModule(module BroadcastServeModule) {
	cluster.serveModulesRW.Lock()
	defer cluster.serveModulesRW.Unlock()
	cluster.serveModules = append(cluster.serveModules, module)
}

func (cluster *stdBroadcastCluster) UnregisterServeModule(module BroadcastServeModule) {
	cluster.serveModulesRW.Lock()
	defer cluster.serveModulesRW.Unlock()
	for i, m := range cluster.serveModules {
		if m == module {
			cluster.serveModules = append(cluster.serveModules[:i], cluster.serveModules[i+1:]...)
			break
		}
	}
}

func (cluster *stdBroadcastCluster) Serve(payload []byte, addr string) error {
	cluster.logger.Debug("BroadcastCluster", "Serve Addr begin: "+addr)
	defer cluster.logger.Debug("BroadcastCluster", "Serve Addr end: "+addr)

	serveModules := cluster.ServeModules()
	if len(serveModules) <= 0 {
		return ErrBroadcastClusterServeUnavailable
	}

	for _, module := range serveModules {
		err := module.ServeBroadcast(payload, addr)
		if err == nil {
			break
		}
	}

	return nil
}

func (cluster *stdBroadcastCluster) DeliverAddrs() []string {
	cluster.deliverAddrsRW.RLock()
	defer cluster.deliverAddrsRW.RUnlock()

	if len(cluster.deliverAddrs) <= 0 {
		return cluster.deliverAddrs
	}

	addrs := make([]string, 0)
	for _, addr := range cluster.deliverAddrs {
		if idx, ok := slices.BinarySearch(addrs, addr); !ok {
			addrs = slices.Insert(addrs, idx, addr)
		}
	}
	return addrs
}

func (cluster *stdBroadcastCluster) SetDeliverAddr(addr []string) {
	cluster.deliverAddrsRW.Lock()
	defer cluster.deliverAddrsRW.Unlock()
	cluster.deliverAddrs = addr
}

func (cluster *stdBroadcastCluster) DeliverLimitSize() int {
	cluster.deliverLimitSizeRW.RLock()
	defer cluster.deliverLimitSizeRW.RUnlock()
	return cluster.deliverLimitSize
}

func (cluster *stdBroadcastCluster) SetDeliverLimitSize(size int) {
	cluster.deliverLimitSizeRW.Lock()
	defer cluster.deliverLimitSizeRW.Unlock()
	cluster.deliverLimitSize = size
}

func (cluster *stdBroadcastCluster) MTU() int {
	cluster.mtuRW.RLock()
	defer cluster.mtuRW.RUnlock()
	return cluster.mtu
}

func (cluster *stdBroadcastCluster) SetMTU(mtu int) {
	cluster.mtuRW.Lock()
	defer cluster.mtuRW.Unlock()
	cluster.mtu = mtu
}

func (cluster *stdBroadcastCluster) Deliver(payload []byte, addrs ...string) error {
	size := len(payload) + 1
	limitSize := cluster.DeliverLimitSize()
	if size > limitSize && limitSize > 0 {
		return bytes.ErrTooLarge
	}

	var deliverAddrs []string
	if len(addrs) <= 0 {
		deliverAddrs = cluster.DeliverAddrs()
	} else {
		deliverAddrs = addrs
	}

	if len(deliverAddrs) <= 0 {
		return ErrBroadcastClusterDeliverNoAddrs
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
	mtu := cluster.MTU()
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
