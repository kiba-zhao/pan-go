// Package peer provides p2p functionality for Pan.
//
// It is used to define p2p functionality in the application
package peer

import (
	"errors"
	"pan/lib/app"
	"pan/lib/bootstrap"
	"pan/lib/config"
	"pan/lib/injection"
	"pan/lib/log"
	"pan/lib/runtime"
	"path/filepath"
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
	AppConfig config.AppConfig

	cluster *stdPeerCluster
	cfg     *stdPeerConfig

	registry runtime.Registry
	rw       sync.RWMutex
	already  bool
}

func New() interface{} {
	module := &stdPeerModule{}

	cluster := &stdPeerCluster{}
	module.cluster = cluster
	cluster.logger = log.Default()
	module.cfg = &stdPeerConfig{}
	return module
}

var _ = (runtime.InitializeModule)((*stdPeerModule)(nil))

// Init initializes the peer module with the provided registry.
//
// It sets the module's registry and does not return an error.
func (pn *stdPeerModule) Init(registry runtime.Registry) error {

	pn.rw.Lock()
	pn.registry = registry
	pn.rw.Unlock()

	if pn.already {
		return pn.ReloadModules()
	}
	return nil
}

var _ = (bootstrap.DeferModule)((*stdPeerModule)(nil))

// Defer reloads the peer modules.
//
// It is called by the runtime to reload the peer modules after the application has finished initializing.
//
// The function first checks if the registry is available, and if it is not, an error is returned.
// If the registry is available, the function calls ReloadModules to reload the peer modules.
//
// ReloadModules is a no-op if the registry is not available.
func (pn *stdPeerModule) Defer() error {

	configPath := filepath.Dir(pn.AppConfig.ConfigFilePath())
	err := pn.cfg.EnsureConfig(configPath)
	if err == nil {
		pn.already = true
		err = pn.ReloadModules()
	}

	return err
}

var _ = (runtime.EngineExtensionModule)((*stdPeerModule)(nil))

// EngineTypes returns a slice of reflect.Type representing the various engine types
// associated with the peer module. These types include:
//
//   - PeerAppModule: Represents a module for the p2p application.

func (pn *stdPeerModule) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[PeerAppModule](),
	}
}

var _ = (injection.ComponentProvider)((*stdPeerModule)(nil))

// Components returns a slice of injection.Component representing the components
// provided by the peer module. Currently, the only component provided is the
// PeerModule itself, which is scoped externally.
func (pn *stdPeerModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(pn, injection.ComponentNoneScope),
		injection.NewComponent[PeerCluster](pn.cluster, injection.ComponentExternalScope),
		injection.NewComponent[PeerConfig](pn.cfg, injection.ComponentExternalScope),
	}
}

func (pn *stdPeerModule) ReloadModules() error {
	pn.rw.RLock()
	registry := pn.registry
	pn.rw.RUnlock()

	peerApp := app.NewApp()

	err := runtime.TraverseRegistry(registry, func(module PeerAppModule) error {
		return module.SetupToPeer(peerApp)
	})

	if err == nil {
		pn.cluster.SetPeerApp(peerApp)
	}
	return err
}
