package peer

import (
	"pan/core"
	"pan/memory"
	"sync"
)

type Provider interface {
	App() core.App[Context]
	DialerManager() DialerManager
	HandshakeManager() HandshakeManager
	Peer() Peer
	Settings() Settings
	PeerIdGenerator() PeerIdGenerator
	PeerEvent() PeerEvent
}

type providerSt struct {
	app              core.App[Context]
	appOnce          sync.Once
	dialerMgr        DialerManager
	dialerMgrOnce    sync.Once
	handshakeMgr     HandshakeManager
	handshakeMgrOnce sync.Once
	peer             Peer
	peerOnce         sync.Once
	generator        PeerIdGenerator
	generatorOnce    sync.Once
	settings         Settings
	settingsOnce     sync.Once
	event            PeerEvent
	eventOnce        sync.Once
}

// NewProvider ...
func NewProvider() Provider {
	provider := new(providerSt)
	return provider
}

func (p *providerSt) App() core.App[Context] {
	p.appOnce.Do(func() {
		p.app = core.New[Context]()
	})
	return p.app
}

func (p *providerSt) DialerManager() DialerManager {
	p.dialerMgrOnce.Do(func() {
		p.dialerMgr = memory.NewMap[uint8, NodeDialer]()
	})
	return p.dialerMgr
}

func (p *providerSt) HandshakeManager() HandshakeManager {
	p.handshakeMgrOnce.Do(func() {
		p.handshakeMgr = memory.NewMap[NodeType, NodeHandshake]()
	})
	return p.handshakeMgr
}

func (p *providerSt) Peer() Peer {
	p.peerOnce.Do(func() {
		p.peer = NewPeer(p)
	})
	return p.peer
}

func (p *providerSt) Settings() Settings {
	p.settingsOnce.Do(func() {
		// TODO: Implement reading from toml file
		p.settings = NewSettings()
	})
	return p.settings
}

func (p *providerSt) PeerIdGenerator() PeerIdGenerator {
	p.generatorOnce.Do(func() {
		generator, err := NewPeerIdGenerator(p)
		if err != nil {
			panic(err)
		}
		p.generator = generator
	})
	return p.generator
}

func (p *providerSt) PeerEvent() PeerEvent {
	p.eventOnce.Do(func() {
		p.event = NewPeerEvent()
	})
	return p.event
}
