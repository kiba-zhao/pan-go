package peer

import (
	"pan/core"
	"pan/memory"
)

const (
	nodeAddedEvent = uint8(iota)
	nodeRemovedEvent
	routeAddedEvent
	routeRemovedEvent
)

type PeerEventHandler interface {
	OnNodeAdded(peerId PeerId)
	OnNodeRemoved(peerId PeerId)
	OnRouteAdded(peerId PeerId)
	OnRouteRemoved(peerId PeerId)
}

type PeerEvent interface {
	PeerEventHandler
	core.Event[PeerEventHandler]
}

type peerEventSt struct {
	pocket *memory.Pocket[PeerEventHandler]
}

// NewPeerEvent ...
func NewPeerEvent() PeerEvent {
	event := new(peerEventSt)
	event.pocket = memory.NewPocket[PeerEventHandler]()
	return event
}

// OnNodeAdded ...
func (e *peerEventSt) OnNodeAdded(peerId PeerId) {
	e.onEvent(nodeAddedEvent, peerId)
}

// OnNodeRemoved ...
func (e *peerEventSt) OnNodeRemoved(peerId PeerId) {
	e.onEvent(nodeRemovedEvent, peerId)
}

// OnRouteAdded ...
func (e *peerEventSt) OnRouteAdded(peerId PeerId) {
	e.onEvent(routeAddedEvent, peerId)
}

// OnRouteRemoved ...
func (e *peerEventSt) OnRouteRemoved(peerId PeerId) {
	e.onEvent(routeRemovedEvent, peerId)
}

// onEvent ...
func (e *peerEventSt) onEvent(name uint8, peerId PeerId) {
	handlers := e.pocket.GetAll()
	for _, handler := range handlers {
		switch name {
		case nodeAddedEvent:
			handler.OnNodeAdded(peerId)
		case nodeRemovedEvent:
			handler.OnNodeRemoved(peerId)
		case routeAddedEvent:
			handler.OnRouteAdded(peerId)
		case routeRemovedEvent:
			handler.OnRouteRemoved(peerId)
		}

	}
}

// Attach ...
func (e *peerEventSt) Attach(handler PeerEventHandler) {
	e.pocket.Add(handler)
}

// Dettach ...
func (e *peerEventSt) Dettach(handler PeerEventHandler) {
	e.pocket.Remove(handler)
}
