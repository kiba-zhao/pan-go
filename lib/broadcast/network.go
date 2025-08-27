package broadcast

import (
	"bytes"
	"errors"
	"net"
	"pan/lib/log"
	"slices"
	"strings"
	"sync"
)

var ErrBroadcastNetworkUnavailable = errors.New("broadcast.BroadcastNetwork Error: Unavailable")
var ErrBroadcastNetworkNoDeliverAddrs = errors.New("broadcast.BroadcastNetwork Error: No Deliver Addrs")
var ErrBroadcastNetworkDeliverInvalidAddrs = errors.New("broadcast.BroadcastNetwork Error: Deliver Invalid Addrs")

// BroadcastServeModule is the broadcast serve module
type BroadcastServeModule interface {
	// ServeBroadcast serves the broadcast message to the peer.
	ServeBroadcast([]byte, string) error
}

type BroadcastNetwork interface {
	// Serve serves the broadcast message to the peer.
	Serve([]byte, string) error

	// Deliver delivers the broadcast message to the peer.
	Deliver([]byte, ...string) error
}

type stdBroadcastNetwork struct {
	logger   log.Logger
	provider *stdBroadcastProvider

	addrs          []string
	deliverMaxSize int
	zoneList       []string
	rw             sync.RWMutex
}

var _ = (BroadcastNetwork)((*stdBroadcastNetwork)(nil))

func (network *stdBroadcastNetwork) Serve(payload []byte, addr string) error {
	network.logger.Debug("BroadcastNetwork", "Serve Addr begin: "+addr)
	defer network.logger.Debug("BroadcastNetwork", "Serve Addr end: "+addr)

	provider := network.provider
	if provider == nil {
		return ErrBroadcastNetworkUnavailable
	}
	serveModules := provider.ServeModules()
	if len(serveModules) <= 0 {
		return ErrBroadcastNetworkUnavailable
	}

	for _, module := range serveModules {
		err := module.ServeBroadcast(payload, addr)
		if err == nil {
			break
		}
	}

	return nil
}

func (network *stdBroadcastNetwork) Deliver(payload []byte, addrs ...string) error {
	network.logger.Debug("BroadcastNetwork", "Deliver Addrs begin:  "+strings.Join(addrs, ", "))
	defer network.logger.Debug("BroadcastNetwork", "Deliver Addrs end")

	provider := network.provider
	if provider == nil {
		return ErrBroadcastNetworkUnavailable
	}

	connSeq := provider.SeqForDeliverConn()
	if connSeq == nil {
		return ErrBroadcastNetworkUnavailable
	}

	size := len(payload) + 1
	network.rw.RLock()
	maxSize := network.deliverMaxSize
	zoneList := network.zoneList
	network.rw.RUnlock()
	if size > maxSize && maxSize > 0 {
		return bytes.ErrTooLarge
	}

	var deliverAddrs []string
	if len(addrs) <= 0 {
		network.rw.RLock()
		deliverAddrs = network.addrs
		network.rw.RUnlock()
	} else {
		deliverAddrs = addrs
	}

	if len(deliverAddrs) <= 0 {
		return ErrBroadcastNetworkNoDeliverAddrs
	}

	deliverUDPAddrs := make([]*net.UDPAddr, 0)
	for _, addr := range deliverAddrs {
		udpAddrs, err := resolveAddrs(addr, zoneList)
		if err != nil {
			network.logger.Error("BroadcastNetwork", "Deliver resolveAddrs error: "+addr)
			continue
		}
		if len(udpAddrs) <= 0 {
			continue
		}
		deliverUDPAddrs = append(deliverUDPAddrs, udpAddrs...)
	}

	if len(deliverUDPAddrs) <= 0 {
		return ErrBroadcastNetworkDeliverInvalidAddrs
	}

	buffer := packBuffer(payload)
	size = len(buffer)

	for mtu, conn := range connSeq {

		localAddr := conn.LocalAddr().(*net.UDPAddr)
		localAddrIP := localAddr.IP

		for _, udpAddr := range deliverUDPAddrs {
			if localAddrIP != nil {
				if localAddrIP.To4() != nil && udpAddr.IP.To4() == nil {
					continue
				}
				if localAddrIP.To4() == nil && udpAddr.IP.To4() != nil {
					continue
				}
			}

			for offset := 0; offset < size; offset += mtu {
				var limit int
				if offset+mtu > size {
					limit = size - offset
				} else {
					limit = offset + mtu
				}
				block := buffer[offset:limit]

				_, err := conn.WriteTo(block, udpAddr)
				if err != nil {
					network.logger.Error("BroadcastNetwork", "Deliver error: "+err.Error()+" Addr: "+udpAddr.String())
				}
			}

		}
	}

	return nil
}

func (network *stdBroadcastNetwork) Setup(config BroadcastConfig) {
	network.logger.Debug("BroadcastNetwork", "Setup")

	network.rw.Lock()
	defer network.rw.Unlock()

	addrs := config.Addrs()
	deliverMaxSize := config.DeliverMaxSize()
	zoneList := config.IPv6ZoneList()

	if !slices.Equal(network.addrs, addrs) {
		network.addrs = addrs
	}

	if network.deliverMaxSize != deliverMaxSize {
		network.deliverMaxSize = deliverMaxSize
	}

	if !slices.Equal(network.zoneList, zoneList) {
		network.zoneList = zoneList
	}
}

func resolveAddrs(addr string, zoneArr []string) ([]*net.UDPAddr, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	if udpAddr.IP.To4() != nil || !udpAddr.IP.IsMulticast() || len(udpAddr.Zone) > 0 {
		return []*net.UDPAddr{udpAddr}, nil
	}

	if len(zoneArr) <= 0 {
		return nil, errors.New("Unsupported IPv6 Zone")
	}

	udpAddrs := make([]*net.UDPAddr, 0)
	for _, zone := range zoneArr {
		udpAddr_ := net.UDPAddrFromAddrPort(udpAddr.AddrPort())
		udpAddr_.Zone = zone
		udpAddrs = append(udpAddrs, udpAddr_)
	}
	return udpAddrs, err
}
