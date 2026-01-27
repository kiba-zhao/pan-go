package quic

import (
	"context"
	"errors"
	"io"
	"iter"
	"net"
	"pan/internal/peer"
	"sync"
)

var ErrQuicExplorerUnavailable = errors.New("quic.QuicExplorer Error: Unavailable")
var ErrQuicExplorerNotFound = errors.New("quic.QuicExplorer Error: Not Found")

type QuicExplorerGuide interface {
	LookupExplorerAddr(peerId peer.PeerID) (iter.Seq[string], error)
}

type QuicExplorer interface {
	Guides() []QuicExplorerGuide
	AddGuide(guide QuicExplorerGuide)
	RemoveGuide(guide QuicExplorerGuide)
}

type stdQuicExplorer struct {
	network  *stdQuicNetwork
	guides   []QuicExplorerGuide
	guidesRW sync.RWMutex
}

var _ = (peer.PeerTransport)((*stdQuicExplorer)(nil))

func (e *stdQuicExplorer) RoundTrip(ctx context.Context, peerId peer.PeerID, reader io.Reader) (io.ReadCloser, error) {

	addrSeq := e.lookup(peerId)
	quicAddrSeq := resolveExplorerAddrs(addrSeq)
	if quicAddrSeq == nil {
		return nil, ErrQuicExplorerUnavailable
	}

	network := e.network
	hasRoute := false
	for addr := range quicAddrSeq {
		err := network.Route(peerId, addr)
		if err == nil {
			hasRoute = true
			break
		}
	}

	if !hasRoute {
		return nil, ErrQuicExplorerNotFound
	}
	return network.RoundTrip(ctx, peerId, reader)
}

func (e *stdQuicExplorer) CanReach(peerId peer.PeerID) bool {
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

var _ = (QuicExplorer)((*stdQuicExplorer)(nil))

func (e *stdQuicExplorer) Guides() []QuicExplorerGuide {
	e.guidesRW.RLock()
	defer e.guidesRW.RUnlock()
	return e.guides
}

func (e *stdQuicExplorer) AddGuide(guide QuicExplorerGuide) {
	e.guidesRW.Lock()
	defer e.guidesRW.Unlock()
	e.guides = append(e.guides, guide)
}

func (e *stdQuicExplorer) RemoveGuide(guide QuicExplorerGuide) {
	e.guidesRW.Lock()
	defer e.guidesRW.Unlock()
	for i, g := range e.guides {
		if g == guide {
			e.guides = append(e.guides[:i], e.guides[i+1:]...)
			break
		}
	}
}

func (e *stdQuicExplorer) lookup(peerId peer.PeerID) iter.Seq[string] {

	guides := e.Guides()
	if len(guides) == 0 {
		return nil
	}

	return func(yield func(string) bool) {
	loop_guide:
		for _, guide := range guides {
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
