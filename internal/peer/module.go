// Package peer provides p2p functionality for Pan.
//
// It is used to define p2p functionality in the application
package peer

import (
	"context"
	"errors"
	"pan/internal/app"
	"pan/internal/bootstrap"
	"pan/internal/injection"
	"pan/internal/log"
	"pan/internal/runtime"
	"reflect"
	"sync"
)

var ErrPeerModuleUnavailable = errors.New("peer.PeerModule Error: Unavailable")
var ErrPeerModuleUnknownPeer = errors.New("peer.PeerModule Error: Unknown Peer")
var ErrPeerModuleControlConflict = errors.New("peer.PeerModule Error: Control Conflict")

// PeerAppModule is a module for the p2p application
//
// The module that implements this interface will be obtained by the application from the runtime and loaded into the p2p engine
type PeerAppModule interface {
	SetupToPeer(PeerRouter) error
}

type stdPeerModule struct {
	network *stdPeerNetwork
	guard   *stdPeerGuard

	registry runtime.Registry
	rw       sync.RWMutex
	already  bool
}

func New() interface{} {
	logger := log.Default()

	network := &stdPeerNetwork{}
	network.logger = logger
	network.transports = make([]PeerTransport, 0)

	guard := &stdPeerGuard{}
	guard.logger = logger
	guard.blackLists = make([]PeerBlackList, 0)

	module := &stdPeerModule{}
	module.network = network
	module.guard = guard

	return module
}

var _ = (runtime.InitializeModule)((*stdPeerModule)(nil))

func (pn *stdPeerModule) Init(ctx context.Context, registry runtime.Registry) error {

	pn.rw.Lock()
	pn.registry = registry
	pn.rw.Unlock()

	if pn.already {
		return pn.ReloadModules(ctx)
	}
	return nil
}

var _ = (bootstrap.DeferModule)((*stdPeerModule)(nil))

func (pn *stdPeerModule) Defer(ctx context.Context) error {

	if pn.already {
		return nil
	}
	pn.already = true
	return pn.ReloadModules(ctx)
}

var _ = (runtime.EngineExtensionModule)((*stdPeerModule)(nil))

func (pn *stdPeerModule) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[PeerAppModule](),
	}
}

var _ = (injection.ComponentProvider)((*stdPeerModule)(nil))

func (pn *stdPeerModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent[PeerNetwork](pn.network, injection.ComponentExternalScope),
		injection.NewComponent[PeerGuard](pn.guard, injection.ComponentExternalScope),
	}
}

func (pn *stdPeerModule) ReloadModules(ctx context.Context) error {
	pn.rw.RLock()
	registry := pn.registry
	pn.rw.RUnlock()

	peerApp := app.NewApp()

	err := runtime.TraverseRegistry(registry, func(module PeerAppModule) error {
		if ctxErr := runtime.EnsureContext(ctx); ctxErr != nil {
			return ctxErr
		}
		return module.SetupToPeer(peerApp)
	})

	if err == nil {
		pn.network.SetupApp(peerApp)
	}
	return err
}
