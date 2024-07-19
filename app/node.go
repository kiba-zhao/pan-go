package app

import (
	"bytes"
	"context"
	"crypto"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"io"
	"os"
	"pan/app/node"
	"pan/runtime"
	"reflect"
	"sync"
)

type NodeApp = *node.App
type NodeRouter = node.AppHandleGroup
type NodeContext = node.AppContext
type NodeNext = node.Next

type NodeAppModule interface {
	SetupToNode(NodeRouter) error
}

type NodeAppModuleProvider interface {
	NodeAppModules() []NodeAppModule
}

type NodeStream interface {
	io.Reader
	io.Writer
	io.Closer
}

type NodeID = []byte
type Node interface {
	ID() NodeID
	Close() error
}

var (
	ContextNode = []byte("NODE")
)

type NodeTransport interface {
	Do(NodeID, io.Reader, context.Context) (io.Reader, error)
}

type NodeLookupTransport interface {
	NodeTransport
	Lookup(NodeID, context.Context) error
}

type NodeGreetTransport interface {
	NodeTransport
	Greet(NodeID, context.Context) error
}

type NodeTripper interface {
	RoundTrip(NodeID, context.Context, []NodeLookupTransport, []NodeGreetTransport) (NodeTransport, error)
}

type nodeTripper struct {
	module *nodeModule
}

func (nt *nodeTripper) RoundTrip(nodeId NodeID, ctx context.Context, lookUpTransports []NodeLookupTransport, greetTransports []NodeGreetTransport) (NodeTransport, error) {

	for _, transport := range lookUpTransports {
		err := transport.Lookup(nodeId, ctx)
		if err == nil {
			return transport, nil
		}
	}

	var wg sync.WaitGroup

	idx := -1
	ctx_, cancel := context.WithCancelCause(ctx)
	greetChan := make(chan int)
	for idx_, transport := range greetTransports {
		select {
		case idx = <-greetChan:
			cancel(nil)
		default:
			wg.Add(1)
			go func(transport NodeGreetTransport, idx_ int) {
				defer wg.Done()
				err := transport.Greet(nodeId, ctx_)
				if err == nil {
					greetChan <- idx_
				}

			}(transport, idx_)
		}
		if idx >= 0 {
			break
		}
	}

	wg.Wait()
	if idx < 0 {
		return nil, ctx_.Err()
	}
	return greetTransports[idx], nil
}

type NodeDoContext struct {
	ctx context.Context
}
type NodeDoContextUpdater = func(NodeDoContext)

func WithNodeDoContext(ctx context.Context) NodeDoContextUpdater {
	return func(doContext NodeDoContext) {
		doContext.ctx = ctx
	}
}

type NodeSettings interface {
	NodeID() NodeID
	PubKey() any
	PrivKey() crypto.PrivateKey
	Certificate() tls.Certificate
	Available() bool
}

type NodeSettingsListener interface {
	OnNodeSettingsUpdated(NodeSettings)
}

type nodeSettings struct {
	registry       runtime.Registry
	registryLocker sync.RWMutex
	locker         sync.RWMutex
	nodeId         NodeID
	pubKey         any
	privKey        crypto.PrivateKey
	cert           tls.Certificate
	hashCode       []byte
	version        uint8
}

func (ns *nodeSettings) Init(registry runtime.Registry) error {
	ns.registryLocker.Lock()
	defer ns.registryLocker.Unlock()
	ns.registry = registry

	ns.locker.RLock()
	defer ns.locker.RUnlock()
	ns.onUpdated()
	return nil
}

func (ns *nodeSettings) OnConfigUpdated(settings AppSettings) {
	ns.locker.Lock()
	defer ns.locker.Unlock()

	var err error
	v := ns.version
	defer func() {
		if err != nil {
			ns.hashCode = nil
		}
		if v == ns.version {
			return
		}

		ns.registryLocker.RLock()
		defer ns.registryLocker.RUnlock()
		if ns.registry != nil {
			ns.onUpdated()
		}
	}()

	keyPEMBlock, err := os.ReadFile(settings.PrivateKeyPath)
	if err != nil {
		ns.version++
		return
	}
	certPEMBlock, err := os.ReadFile(settings.CertificatePath)
	if err != nil {
		ns.version++
		return
	}

	hash := sha512.New()
	hash.Write(certPEMBlock)
	hash.Write(keyPEMBlock)
	hashCode := hash.Sum(nil)
	if ns.hashCode != nil && bytes.Equal(ns.hashCode, hashCode) {
		return
	}

	ns.version++
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return
	}
	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return
	}
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(x509Cert.PublicKey)
	if err == nil {
		ns.nodeId = pubKeyBytes
		ns.pubKey = x509Cert.PublicKey
		ns.privKey = cert.PrivateKey
		ns.cert = cert
		ns.hashCode = hashCode
	}

	return
}

func (ns *nodeSettings) onUpdated() {
	listeners := runtime.ModulesForType[NodeSettingsListener](ns.registry)
	for _, listener := range listeners {
		listener.OnNodeSettingsUpdated(ns)
	}
}

func (ns *nodeSettings) NodeID() NodeID {
	ns.locker.RLock()
	defer ns.locker.RUnlock()
	return ns.nodeId
}

func (ns *nodeSettings) PubKey() any {
	ns.locker.RLock()
	defer ns.locker.RUnlock()
	return ns.pubKey
}

func (ns *nodeSettings) PrivKey() crypto.PrivateKey {
	ns.locker.RLock()
	defer ns.locker.RUnlock()
	return ns.privKey
}

func (ns *nodeSettings) Certificate() tls.Certificate {

	ns.locker.RLock()
	defer ns.locker.RUnlock()
	return ns.cert
}

func (ns *nodeSettings) Available() bool {
	ns.locker.RLock()
	defer ns.locker.RUnlock()
	return ns.hashCode != nil
}

func (ns *nodeSettings) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[NodeSettingsListener](),
	}
}

type NodeModule interface {
	Serve(NodeStream, Node) error
	Do(NodeID, *node.Request, ...NodeDoContextUpdater) (*node.Response, error)
	Request(NodeID, node.RequestName, io.Reader, ...NodeDoContextUpdater) (*node.Response, error)
	NodeTripper() NodeTripper
	SetNodeTripper(NodeTripper)
	NodeSettings() NodeSettings
	ReloadModules() error
}

type nodeModule struct {
	settings           NodeSettings
	settingsOnce       sync.Once
	tripper            NodeTripper
	tripperLocker      sync.RWMutex
	registry           runtime.Registry
	registryLocker     sync.RWMutex
	defaultTripper     NodeTripper
	defaultTripperOnce sync.Once
	app                NodeApp
	appLocker          sync.RWMutex
}

func (nm *nodeModule) DefaultNodeTripper() NodeTripper {

	nm.defaultTripperOnce.Do(func() {
		nm.tripperLocker.Lock()
		defer nm.tripperLocker.Unlock()
		nm.defaultTripper = &nodeTripper{module: nm}
	})

	nm.tripperLocker.RLock()
	defer nm.tripperLocker.RUnlock()
	return nm.defaultTripper
}

func (nm *nodeModule) Init(registry runtime.Registry) error {

	nm.registryLocker.Lock()
	nm.registry = registry
	nm.registryLocker.Unlock()

	return nm.ReloadModules()
}

func (nm *nodeModule) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[NodeAppModule](),
		reflect.TypeFor[NodeAppModuleProvider](),
		reflect.TypeFor[NodeLookupTransport](),
		reflect.TypeFor[NodeGreetTransport](),
	}
}

func (nm *nodeModule) Components() []runtime.Component {
	return []runtime.Component{
		runtime.NewComponent[NodeModule](nm, runtime.ComponentExternalScope),
	}
}

func (nm *nodeModule) Modules() []interface{} {
	return []interface{}{
		nm.NodeSettings(),
	}
}

func (nm *nodeModule) NodeTripper() NodeTripper {
	nm.tripperLocker.RLock()
	if nm.tripper != nil {
		defer nm.tripperLocker.RUnlock()
		return nm.tripper
	}
	nm.tripperLocker.RUnlock()
	return nm.DefaultNodeTripper()
}

func (nm *nodeModule) SetNodeTripper(tripper NodeTripper) {
	nm.tripperLocker.Lock()
	defer nm.tripperLocker.Unlock()
	nm.tripper = tripper
}

func (nm *nodeModule) Serve(stream NodeStream, target Node) error {

	var app NodeApp
	nm.appLocker.RLock()
	if nm.app != nil {
		app = nm.app
	}
	nm.appLocker.RUnlock()

	ctx := node.NewAppContext()
	ctx.Set(ContextNode, target)
	err := node.UnmarshalRequest(stream, ctx.Request())
	if err == nil {
		if app == nil {
			err = ErrUnavailable
		} else {
			err = app.Run(ctx, nil)
		}
	}

	if err != nil {
		ctx.ThrowError(CodeInternalError, err)
	}

	if ctx.Code() < 0 {
		ctx.ThrowError(CodeNotFound, ErrNotFound)
	}

	reader := node.MarshalResponse(&ctx.Response)
	_, resErr := io.Copy(stream, reader)
	if resErr == nil {
		resErr = stream.Close()
	}

	if err == nil && resErr != nil {
		err = resErr
	}
	return err
}

func (nm *nodeModule) Do(nodeId NodeID, request *node.Request, updaters ...NodeDoContextUpdater) (*node.Response, error) {

	doContext := NodeDoContext{}
	for _, updater := range updaters {
		updater(doContext)
	}
	if doContext.ctx == nil {
		doContext.ctx = context.Background()
	}

	nm.registryLocker.RLock()
	lookUpTransports := runtime.ModulesForType[NodeLookupTransport](nm.registry)
	greetTransports := runtime.ModulesForType[NodeGreetTransport](nm.registry)
	nm.registryLocker.RUnlock()

	ctx := doContext.ctx
	tripper := nm.NodeTripper()
	reqReader := node.MarshalRequest(request)

	var err error
	var transport NodeTransport
	var resReader io.Reader
	for {
		transport, err = tripper.RoundTrip(nodeId, ctx, lookUpTransports, greetTransports)
		if err != nil {
			break
		}

		resReader, err = transport.Do(nodeId, reqReader, ctx)
		if err == nil || ctx.Err() != nil {
			break
		}

	}

	if err != nil {
		return nil, err
	}

	response := &node.Response{}
	node.InitResponse(response)
	err = node.UnmarshalResponse(resReader, response)

	return response, err
}

func (nm *nodeModule) Request(nodeId NodeID, name node.RequestName, body io.Reader, updaters ...NodeDoContextUpdater) (*node.Response, error) {
	request := node.NewRequest(name, body)
	return nm.Do(nodeId, request, updaters...)
}

func (nm *nodeModule) NodeSettings() NodeSettings {
	nm.settingsOnce.Do(func() {
		nm.settings = &nodeSettings{}
	})
	return nm.settings
}

func (nm *nodeModule) ReloadModules() error {
	nm.registryLocker.RLock()
	registry := nm.registry
	nm.registryLocker.RUnlock()

	nm.appLocker.Lock()
	defer nm.appLocker.Unlock()
	app := node.NewApp()

	err := runtime.TraverseRegistry(registry, func(module NodeAppModule) error {
		return module.SetupToNode(nm.app)
	})
	if err != nil {
		return err
	}

	err = runtime.TraverseRegistry(registry, func(module NodeAppModuleProvider) error {
		var setupErr error
		for _, m := range module.NodeAppModules() {
			setupErr = m.SetupToNode(nm.app)
			if err != nil {
				break
			}
		}
		return setupErr
	})

	if err == nil {
		nm.app = app
	}
	return err
}
