package quic

import (
	"errors"
	"iter"
	"maps"
	"net"
	"sync"

	"github.com/quic-go/quic-go"
)

type stdQuicTransportProvider struct {
	transportMap map[string]*QuicTransport
	transportsRW sync.RWMutex
}

func (provider *stdQuicTransportProvider) SeqForTransports() iter.Seq2[string, *QuicTransport] {
	provider.transportsRW.RLock()
	defer provider.transportsRW.RUnlock()

	transportMap := maps.Clone(provider.transportMap)
	return func(yield func(string, *QuicTransport) bool) {
		for addr, transport := range transportMap {
			if !yield(addr, transport) {
				break
			}
		}
	}
}

func (provider *stdQuicTransportProvider) NewTransport(addr string, port uint16) (*QuicTransport, error) {
	provider.transportsRW.Lock()
	defer provider.transportsRW.Unlock()

	var ipNet *net.IPNet
	var addrIP net.IP
	var err error
	if len(addr) > 0 {
		if _, ok := provider.transportMap[addr]; ok {
			return nil, errors.New("quic.QuicRuntime Error: Transport Already Exists On " + addr)
		}

		addrIP, ipNet, err = net.ParseCIDR(addr)
		if err != nil {
			return nil, err
		}
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   addrIP,
		Port: int(port),
	})
	if err != nil {
		return nil, err
	}

	transport := &QuicTransport{}
	transport.ipNet = ipNet
	transport.Transport = &quic.Transport{
		Conn: conn,
	}

	provider.transportMap[addr] = transport
	return transport, nil
}

func (provider *stdQuicTransportProvider) RevokeTransport(addr string) error {
	provider.transportsRW.Lock()
	defer provider.transportsRW.Unlock()
	transport, ok := provider.transportMap[addr]
	if !ok {
		return nil
	}
	delete(provider.transportMap, addr)
	return transport.Close()
}

func (provider *stdQuicTransportProvider) RevokeAllTransports() error {
	provider.transportsRW.Lock()
	defer provider.transportsRW.Unlock()

	transportMap := provider.transportMap
	provider.transportMap = make(map[string]*QuicTransport)

	errs := make([]error, 0)
	for _, transport := range transportMap {
		err := transport.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) <= 0 {
		return nil
	}
	return errors.Join(errs...)
}

func (provider *stdQuicTransportProvider) SeqForTransportsWithAddrIP(addrIP net.IP) iter.Seq2[string, *QuicTransport] {
	provider.transportsRW.RLock()
	defer provider.transportsRW.RUnlock()

	transportMap := maps.Clone(provider.transportMap)

	return func(yield func(string, *QuicTransport) bool) {

		for addr, transport := range transportMap {
			ipNet := transport.ipNet
			if ipNet != nil {
				if (addrIP.To4() == nil && ipNet.IP.To4() != nil) || (addrIP.To4() != nil && ipNet.IP.To4() == nil) {
					continue
				}
				if !addrIP.IsGlobalUnicast() && !ipNet.Contains(addrIP) {
					continue
				}
				if addrIP.IsGlobalUnicast() && ipNet.IP.IsLoopback() {
					continue
				}
			}

			if !yield(addr, transport) {
				break
			}
		}
	}
}
