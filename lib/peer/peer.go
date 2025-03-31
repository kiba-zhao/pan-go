// Package peer provides p2p functionality for Pan.
//
// It is used to define p2p functionality in the application
package peer

import (
	"context"
	"errors"
	"io"
	"pan/lib/injection"
	"pan/lib/runtime"
	"reflect"
	"sync"
)

const (
	CodeOK            = 0
	CodeInternalError = 500
	CodeNotFound      = 404
	CodeBadRequest    = 400
	CodeForbidden     = 403
)

var ErrPeerModuleUnavailable = errors.New("peer.PeerModule Error: Unavailable")
var ErrPeerModuleInvalidApp = errors.New("peer.PeerModule Error: Invalid App")
var ErrPeerModuleNotFound = errors.New("peer.PeerModule Error: Serve Not Found")
var ErrPeerModuleUnknownPeer = errors.New("peer.PeerModule Error: Unknown Peer")
var ErrPeerModuleControlConflict = errors.New("peer.PeerModule Error: Control Conflict")
var ErrPeerModuleNetworkNotFound = errors.New("peer.PeerModule Error: Network Not Found")
var ErrPeerModuleAccessDenied = errors.New("peer.PeerModule Error: Access Denied")

type PeerID = []byte
type PeerApp = *App
type PeerRouter = AppHandleGroup
type PeerContext = AppContext
type PeerNext = Next

// PeerAppModule is a module for the p2p application
//
// The module that implements this interface will be obtained by the application from the runtime and loaded into the p2p engine
type PeerAppModule interface {
	SetupToPeer(PeerRouter) error
}

// PeerScopeModule is a module that provides a scope for the p2p application
type PeerScopeModule interface {
	PeerScope() []byte
}

// PeerAppModuleProvider is a module that provides multiple PeerAppModules
type PeerAppModuleProvider interface {
	PeerAppModules() []PeerAppModule
}

type PeerStream interface {
	io.Reader
	io.Writer
	io.Closer
}

var (
	ContextPeerID = []byte("PeerID")
)

// PeerGuard is a module that provides access control and security features within the p2p module
type PeerGuard interface {
	Enabled() bool
	Access(PeerID) error
}

// PeerModule is the main interface for the p2p module
type PeerModule interface {
	CanReach(PeerID) bool
	Purge(PeerID) error
	Serve(PeerStream, PeerID) error
	Do(context.Context, PeerID, *Request) (*Response, error)
	Request(context.Context, PeerID, RequestName, io.Reader, ...HeaderItem) (*Response, error)
	PeerSettings() PeerSettings
	ReloadModules() error
	Access(PeerID) error
}

// New creates a new p2p application module.
//
// The module implements the PeerModule interface and is used to manage p2p connections.
//
// The module is also a runtime.Module, and can be used to load components into the p2p engine.
func New() interface{} {
	module := &peerModule{}
	network := &peerNetwork{peerModule: module}
	module.network = network
	return module
}

type peerModule struct {
	settings       PeerSettings
	settingsOnce   sync.Once
	registry       runtime.Registry
	registryLocker sync.RWMutex
	app            PeerApp
	appLocker      sync.RWMutex
	network        *peerNetwork
}

// Init initializes the peer module with the provided registry.
//
// It sets the module's registry and does not return an error.
func (pn *peerModule) Init(registry runtime.Registry) error {

	pn.registryLocker.Lock()
	pn.registry = registry
	pn.registryLocker.Unlock()

	return nil
}

// Defer reloads the peer modules.
//
// It is called by the runtime to reload the peer modules after the application has finished initializing.
//
// The function first checks if the registry is available, and if it is not, an error is returned.
// If the registry is available, the function calls ReloadModules to reload the peer modules.
//
// ReloadModules is a no-op if the registry is not available.
func (pn *peerModule) Defer() error {
	pn.registryLocker.RLock()
	registry := pn.registry
	pn.registryLocker.RUnlock()
	if registry == nil {
		return ErrPeerModuleUnavailable
	}

	return pn.ReloadModules()
}

// EngineTypes returns a slice of reflect.Type representing the various engine types
// associated with the peer module. These types include:
//
//   - PeerNetwork: Represents the network aspect of the peer module.
//   - PeerAppModule: Represents a module for the p2p application.
//   - PeerAppModuleProvider: Provides multiple PeerAppModules.
//   - PeerGuard: Provides access control and security features within the peer module.

func (pn *peerModule) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[PeerNetwork](),
		reflect.TypeFor[PeerAppModule](),
		reflect.TypeFor[PeerAppModuleProvider](),
		reflect.TypeFor[PeerGuard](),
	}
}

// Components returns a slice of injection.Component representing the components
// provided by the peer module. Currently, the only component provided is the
// PeerModule itself, which is scoped externally.
func (pn *peerModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent[PeerModule](pn, injection.ComponentExternalScope),
	}
}

// Modules returns sub-modules of the peer module.
func (pn *peerModule) Modules() []interface{} {
	return []interface{}{
		pn.PeerSettings(),
	}
}

// Serve handles p2p streams of connections
func (pn *peerModule) Serve(stream PeerStream, target PeerID) error {

	defer stream.Close()

	var app PeerApp
	pn.appLocker.RLock()
	if pn.app != nil {
		app = pn.app
	}
	pn.appLocker.RUnlock()

	var err error
	ctx := NewAppContext()
	if app == nil {
		err = ErrPeerModuleInvalidApp
	} else {
		err = UnmarshalRequest(stream, ctx.Request())
	}

	if err == nil {
		ctx.Set(ContextPeerID, target)
		err = app.Run(ctx, nil)
		defer ctx.Close()
	}

	if err != nil {
		ctx.ThrowError(CodeInternalError, err)
	}

	if ctx.Code() < 0 {
		ctx.ThrowError(CodeNotFound, ErrPeerModuleNotFound)
	}

	reader := MarshalResponse(&ctx.Response)
	_, resErr := io.Copy(stream, reader)

	if err == nil && resErr != nil {
		err = resErr
	}

	return err
}

// CanReach checks if a peer is reachable
func (pn *peerModule) CanReach(peerId PeerID) bool {
	return pn.network.CanReach(peerId)
}

// Purge purges a peer with the given peerId
func (pn *peerModule) Purge(peerId PeerID) error {
	return pn.network.Purge(peerId)
}

// Do sends a request to a peer and returns the response
func (pn *peerModule) Do(ctx context.Context, peerId PeerID, request *Request) (*Response, error) {

	reqReader := MarshalRequest(request)
	resReader, err := pn.network.RoundTrip(ctx, peerId, reqReader)

	if err != nil {
		return nil, err
	}

	if resReader == nil {
		return nil, ErrPeerModuleNetworkNotFound
	}

	response := &Response{}
	InitResponse(response)
	err = UnmarshalResponse(resReader, response)

	if err == nil && response.Code() != CodeOK {
		var content []byte
		content, err = io.ReadAll(response)
		if err == nil {
			err = &PeerError{code: response.Code(), err: string(content)}
		}
	}

	if err != nil {
		resReader.Close()
		return nil, err
	}

	return response, err
}

// Request sends a request to a peer and returns the response
func (pn *peerModule) Request(ctx context.Context, peerId PeerID, name RequestName, body io.Reader, headerItems ...HeaderItem) (*Response, error) {
	request := NewRequest(name, body)
	if len(headerItems) > 0 {
		for _, headerItem := range headerItems {
			request.SetHeader(headerItem.Key, headerItem.Value)
		}
	}
	return pn.Do(ctx, peerId, request)
}

// PeerNetworks returns a slice of PeerNetwork modules loaded into the p2p engine
func (pn *peerModule) PeerNetworks() []interface{} {
	pn.registryLocker.RLock()
	registry := pn.registry
	pn.registryLocker.RUnlock()

	if registry == nil {
		return nil
	}

	modules, _ := registry.ModulesByType(reflect.TypeFor[PeerNetwork]())
	return modules
}

// PeerSettings returns the current peer settings.
//
// The settings are initialized on first call and never changed.
func (pn *peerModule) PeerSettings() PeerSettings {
	pn.settingsOnce.Do(func() {
		pn.settings = &peerSettings{}
	})
	return pn.settings
}

// ReloadModules reloads the peer application modules from the registry.
//
// It traverses the registry, looking for modules that implement the PeerAppModule,
// PeerAppModuleProvider, and PeerScopeModule interfaces. For each module,
// it calls the SetupToPeer method to set up the peer application.
//
// If any error occurs during the reloading process, it will be returned.
//
// If the reloading process is successful, the peer application will be updated
// with the new modules.
func (pn *peerModule) ReloadModules() error {
	pn.registryLocker.RLock()
	registry := pn.registry
	pn.registryLocker.RUnlock()

	pn.appLocker.Lock()
	defer pn.appLocker.Unlock()
	app := NewApp()

	err := runtime.TraverseRegistry(registry, func(module PeerAppModule) error {
		return module.SetupToPeer(app)
	})
	if err != nil {
		return err
	}

	err = runtime.TraverseRegistry(registry, func(module PeerAppModuleProvider) error {
		var setupErr error
		var router PeerRouter
		if scopeModule, ok := module.(PeerScopeModule); ok {
			scope := scopeModule.PeerScope()
			router = app.Route(scope)
		} else {
			router = app.Route(nil)
		}

		for _, m := range module.PeerAppModules() {
			setupErr = m.SetupToPeer(router)
			if err != nil {
				break
			}
		}
		return setupErr
	})

	if err == nil {
		pn.app = app
	}
	return err
}

// Access checks if a peer is allowed to access the p2p network.
//
// It traverses the registry, looking for modules that implement the PeerGuard interface.
// For each module, it checks if the module is enabled and if the peer is allowed to
// access the p2p network. If any error occurs during the checking process, it will be
// returned.
//
// If the checking process is successful, the peer is allowed to access the p2p network.
func (pn *peerModule) Access(peerId PeerID) error {
	pn.registryLocker.RLock()
	defer pn.registryLocker.RUnlock()

	return runtime.TraverseRegistry(pn.registry, func(module PeerGuard) error {
		if !module.Enabled() {
			return nil
		}
		return module.Access(peerId)
	})
}
