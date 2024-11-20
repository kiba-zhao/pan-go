package peer

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"pan/app/bootstrap"
	"pan/runtime"
	"reflect"
	"sync"
	"time"
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
	ContextNode = []byte("NODE")
)

type PeerDoContext struct {
	ctx context.Context
}
type PeerDoContextUpdater = func(PeerDoContext)

func WithPeerDoContext(ctx context.Context) PeerDoContextUpdater {
	return func(doContext PeerDoContext) {
		doContext.ctx = ctx
	}
}

type PeerGuard interface {
	Enabled() bool
	Access(PeerID) error
}

type PeerModule interface {
	Serve(PeerStream, PeerNode) error
	Do(PeerID, *Request, ...PeerDoContextUpdater) (*Response, error)
	Request(PeerID, RequestName, io.Reader, ...PeerDoContextUpdater) (*Response, error)
	PeerSettings() PeerSettings
	PeerManager() PeerManager
	ReloadModules() error
	Access(PeerID) error
	Control(PeerNode) error
	NewResourceID(PeerType) PeerResourceID
}

type peerModule struct {
	mgr            PeerManager
	mgrOnce        sync.Once
	settings       PeerSettings
	settingsOnce   sync.Once
	registry       runtime.Registry
	registryLocker sync.RWMutex
	app            PeerApp
	appLocker      sync.RWMutex
	seq            uint32
	seqLocker      sync.Mutex
}

func New() interface{} {
	return &peerModule{}
}

func (pn *peerModule) PeerManager() PeerManager {
	pn.mgrOnce.Do(func() {
		pn.mgr = &peerManager{}
	})
	return pn.mgr
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
		pn.PeerManager(),
	}
}

func (pn *peerModule) Serve(stream PeerStream, target PeerNode) error {

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
		ctx.Set(ContextNode, target)
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

func (pn *peerModule) Do(peerId PeerID, request *Request, updaters ...PeerDoContextUpdater) (*Response, error) {

	doContext := PeerDoContext{}
	for _, updater := range updaters {
		updater(doContext)
	}
	if doContext.ctx == nil {
		doContext.ctx = context.Background()
	}

	ctx := doContext.ctx
	reqReader := MarshalRequest(request)

	resReader, err := pn.roundTrip(ctx, peerId, reqReader)

	if err != nil {
		return nil, err
	}

	response := &Response{}
	InitResponse(response)
	err = UnmarshalResponse(resReader, response)

	return response, err
}

func (pn *peerModule) Request(peerId PeerID, name RequestName, body io.Reader, updaters ...PeerDoContextUpdater) (*Response, error) {
	request := NewRequest(name, body)
	return pn.Do(peerId, request, updaters...)
}

func (pn *peerModule) roundTrip(ctx context.Context, peerId PeerID, reqReader io.Reader) (reader io.ReadCloser, err error) {

	mgr := pn.PeerManager()
	mgr.TraversePeerNode(peerId, func(peerNode PeerNode) bool {
		if !peerNode.IsIdle() {
			return true
		}
		reader, err = peerNode.Do(ctx, reqReader)
		if err != nil && !errors.Is(err, ctx.Err()) {
			return true
		}
		if err != nil {
			peerNode.Close()
		}
		return false
	})

	if err == nil && reader == nil {
		err = ErrPeerModuleUnknownPeer
	}
	return
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

func (pn *peerModule) Control(peerNode PeerNode) error {
	err := pn.Access(peerNode.PeerID())

	if err == nil {
		mgr := pn.PeerManager()
		_, ok := mgr.SearchOrStore(peerNode)
		if ok {
			err = ErrPeerModuleControlConflict
		}
	}

	return err
}

func (pn *peerModule) NewResourceID(peerType PeerType) PeerResourceID {

	pn.seqLocker.Lock()
	pn.seq++
	seq := pn.seq
	pn.seqLocker.Unlock()

	resourceId := make([]byte, 13)
	resourceId[0] = byte(peerType)
	binary.BigEndian.PutUint64(resourceId[1:], uint64(time.Now().Unix()))
	binary.BigEndian.PutUint32(resourceId[9:], seq)

	return PeerResourceID(resourceId)
}
