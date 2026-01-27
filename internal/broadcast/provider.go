package broadcast

import (
	"iter"
	"net"
	"pan/internal/quic"
	"sync"
)

type stdBroadcastProvider struct {
	quicNetwork   quic.QuicNetwork
	quicNetworkRW sync.RWMutex

	ifaces []BroadcastInterface
	rw     sync.RWMutex

	serveModules   []BroadcastServeModule
	serveModulesRW sync.RWMutex
}

func (p *stdBroadcastProvider) SeqForDeliverConn() iter.Seq2[int, *net.UDPConn] {
	p.rw.Lock()
	defer p.rw.Unlock()
	if len(p.ifaces) <= 0 || p.quicNetwork == nil {
		return nil
	}

	ifaces := p.ifaces
	quicNetwork := p.quicNetwork
	return func(yield func(int, *net.UDPConn) bool) {
		for _, iface := range ifaces {
			transport := quicNetwork.SelectTransport(iface.Addr)
			if transport == nil {
				continue
			}
			conn, ok := transport.Conn.(*net.UDPConn)
			if !ok {
				continue
			}
			if !yield(iface.MTU, conn) {
				break
			}
		}
	}
}

func (p *stdBroadcastProvider) Setup(config BroadcastConfig) error {
	p.rw.Lock()
	defer p.rw.Unlock()

	p.ifaces = config.Interfaces()
	return nil
}

func (p *stdBroadcastProvider) ServeModules() []BroadcastServeModule {
	p.serveModulesRW.RLock()
	defer p.serveModulesRW.RUnlock()
	return p.serveModules
}

func (p *stdBroadcastProvider) RegisterServeModule(module BroadcastServeModule) {
	p.serveModulesRW.Lock()
	defer p.serveModulesRW.Unlock()
	p.serveModules = append(p.serveModules, module)
}

func (p *stdBroadcastProvider) UnregisterServeModule(module BroadcastServeModule) {
	p.serveModulesRW.Lock()
	defer p.serveModulesRW.Unlock()
	for i, m := range p.serveModules {
		if m == module {
			p.serveModules = append(p.serveModules[:i], p.serveModules[i+1:]...)
			break
		}
	}
}

func (p *stdBroadcastProvider) SetupQuicNetwork(network quic.QuicNetwork) {
	p.quicNetworkRW.Lock()
	defer p.quicNetworkRW.Unlock()
	p.quicNetwork = network
}

func (p *stdBroadcastProvider) QuicNetwork() quic.QuicNetwork {
	p.quicNetworkRW.RLock()
	defer p.quicNetworkRW.RUnlock()
	return p.quicNetwork
}
