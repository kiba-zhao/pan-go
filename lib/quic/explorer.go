package quic

import (
	"context"
	"errors"
	"io"
	"iter"
	"net"
	"pan/lib/peer"
	"pan/lib/runtime"
	reflect "reflect"
	"sync"
)

var ErrQuicExplorerUnavailable = errors.New("quic.QuicExplorer Error: Unavailable")
var ErrQuicExplorerNotFound = errors.New("quic.QuicExplorer Error: Not Found")

type QuicExplorerGuide interface {
	LookupExplorerAddr(peerId peer.PeerID) (iter.Seq[string], error)
}

type QuicExplorer struct {
	quicModule *quicPeerModule
	registry   runtime.Registry
	rw         sync.RWMutex
}

func (e *QuicExplorer) Init(registry runtime.Registry) error {
	e.rw.Lock()
	e.registry = registry
	e.rw.Unlock()
	return nil
}

func (e *QuicExplorer) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[QuicExplorerGuide](),
	}
}

func (e *QuicExplorer) lookup(peerId peer.PeerID) iter.Seq[string] {
	e.rw.RLock()
	registry := e.registry
	e.rw.RUnlock()

	if registry == nil {
		return nil
	}

	guideSeq := runtime.SeqForType[QuicExplorerGuide](e.registry)

	return func(yield func(string) bool) {
	loop_guide:
		for guide := range guideSeq {
			// discover addr from guide
			addrs, err := guide.LookupExplorerAddr(peerId)
			if err != nil || addrs == nil {
				continue
			}
			//

			for addr := range addrs {
				if !yield(addr) {
					break loop_guide
				}
			}

		}
	}
}

func (e *QuicExplorer) RoundTrip(ctx context.Context, peerId peer.PeerID, reader io.Reader) (io.ReadCloser, error) {

	addrSeq := e.lookup(peerId)
	quicAddrSeq := resolveExplorerAddrs(addrSeq)
	if quicAddrSeq == nil {
		return nil, ErrQuicExplorerUnavailable
	}

	quicModule := e.quicModule
	hasRoute := false
	for addr := range quicAddrSeq {
		err := quicModule.Route(peerId, addr)
		if err == nil {
			hasRoute = true
			break
		}
	}

	if !hasRoute {
		return nil, ErrQuicExplorerNotFound
	}
	return quicModule.RoundTrip(ctx, peerId, reader)
}

func (e *QuicExplorer) CanReach(peerId peer.PeerID) bool {
	addrSeq := e.lookup(peerId)
	if addrSeq == nil {
		return false
	}
	reachable := false
	for _ = range addrSeq {
		reachable = true
		break
	}
	return reachable
}

func resolveExplorerAddrs(addrs iter.Seq[string]) iter.Seq[string] {
	if addrs == nil {
		return nil
	}
	return func(yield func(string) bool) {
	addrs_loop:
		for addr := range addrs {
			var ipAddrs []string
			host, port, err := net.SplitHostPort(addr)
			if addrErr, ok := err.(*net.AddrError); ok {
				switch addrErr.Err {
				case "missing port in address": // resolve host name to ip address
					ipAddrs, err = net.LookupHost(addr)
				case "too many colons in address": // resolve ip address
					_, err = net.ResolveIPAddr("ip", addr)
					if err == nil {
						ipAddrs = []string{addr}
					}
				default:
				}
			}
			if err != nil {
				continue
			}
			if len(ipAddrs) <= 0 {
				ipAddrs = []string{host}
			}
			if len(port) <= 0 {
				port = "9000"
			}

			for _, ipAddr := range ipAddrs {
				quicAddr := net.JoinHostPort(ipAddr, port)
				if !yield(quicAddr) {
					break addrs_loop
				}
			}
		}
	}
}
