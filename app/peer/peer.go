package peer

import (
	"context"
	"errors"
	"io"
	"pan/app/bootstrap"
	"pan/runtime"
	"reflect"
	"sync"
)

const (
	CodeOK            = 0
	CodeInternalError = 500
	CodeNotFound      = 404
	CodeBadRequest    = 400
)

var ErrPeerModuleUnavailable = errors.New("peer.PeerModule Error: Unavailable")
var ErrPeerModuleInvalidApp = errors.New("peer.PeerModule Error: Invalid App")
var ErrPeerModuleNotFound = errors.New("peer.PeerModule Error: Serve Not Found")
var ErrPeerModuleUnknownPeer = errors.New("peer.PeerModule Error: Unknown Peer")
var ErrPeerModuleControlConflict = errors.New("peer.PeerModule Error: Control Conflict")
var ErrPeerModuleNetworkNotFound = errors.New("peer.PeerModule Error: Network Not Found")

type PeerID = []byte
type PeerApp = *App
type PeerRouter = AppHandleGroup
type PeerContext = AppContext
type PeerNext = Next

type PeerAppModule interface {
	SetupToPeer(PeerRouter) error
}

type PeerScopeModule interface {
	PeerScope() []byte
}

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

type PeerGuard interface {
	Enabled() bool
	Access(PeerID) error
}

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
	network        PeerNetwork
}

func (pn *peerModule) Init(registry runtime.Registry) error {

	pn.registryLocker.Lock()
	pn.registry = registry
	pn.registryLocker.Unlock()

	return nil
}

func (pn *peerModule) Defer() error {
	pn.registryLocker.RLock()
	registry := pn.registry
	pn.registryLocker.RUnlock()
	if registry == nil {
		return ErrPeerModuleUnavailable
	}

	return pn.ReloadModules()
}

func (pn *peerModule) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[PeerNetwork](),
		reflect.TypeFor[PeerAppModule](),
		reflect.TypeFor[PeerAppModuleProvider](),
		reflect.TypeFor[PeerGuard](),
	}
}

func (pn *peerModule) Components() []bootstrap.Component {
	return []bootstrap.Component{
		bootstrap.NewComponent[PeerModule](pn, bootstrap.ComponentExternalScope),
	}
}

func (pn *peerModule) Modules() []interface{} {
	return []interface{}{
		pn.PeerSettings(),
	}
}

func (pn *peerModule) Serve(stream PeerStream, target PeerID) error {

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
	if resErr == nil {
		resErr = stream.Close()
	}

	if err == nil && resErr != nil {
		err = resErr
	}

	return err
}

func (pn *peerModule) CanReach(peerId PeerID) bool {
	return pn.network.CanReach(peerId)
}

func (pn *peerModule) Purge(peerId PeerID) error {
	return pn.network.Purge(peerId)
}

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
			err = errors.New(string(content))
		}
	}

	if err != nil {
		resReader.Close()
		return nil, err
	}

	return response, err
}

func (pn *peerModule) Request(ctx context.Context, peerId PeerID, name RequestName, body io.Reader, headerItems ...HeaderItem) (*Response, error) {
	request := NewRequest(name, body)
	if len(headerItems) > 0 {
		for _, headerItem := range headerItems {
			request.SetHeader(headerItem.Key, headerItem.Value)
		}
	}
	return pn.Do(ctx, peerId, request)
}

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

func (pn *peerModule) PeerSettings() PeerSettings {
	pn.settingsOnce.Do(func() {
		pn.settings = &peerSettings{}
	})
	return pn.settings
}

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
