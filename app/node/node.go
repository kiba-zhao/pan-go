package node

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"io"
	"math/big"
	"os"
	"pan/app/bootstrap"
	"pan/app/config"
	"pan/app/constant"
	"pan/runtime"
	"path"
	"reflect"
	"slices"
	"sync"
	"time"
)

type NodeApp = *App
type NodeRouter = AppHandleGroup
type NodeContext = AppContext
type NodeNext = Next

type NodeAppModule interface {
	SetupToNode(NodeRouter) error
}

type NodeScopeModule interface {
	NodeScope() []byte
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
type NodeType = uint8
type NodeResourceID = []byte

const (
	NodeTypeAlive NodeType = iota + 1
	NodeTypeReachable
)

type Node interface {
	ID() NodeID
	Type() NodeType
	Do(context.Context, io.Reader) (io.Reader, error)
	Greet(context.Context) error
	Close() error
	ResourceID() NodeResourceID
}

var (
	ContextNode = []byte("NODE")
)

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
	ns.registry = registry
	ns.registryLocker.Unlock()

	ns.locker.RLock()
	defer ns.locker.RUnlock()
	ns.onUpdated()
	return nil
}

func (ns *nodeSettings) OnConfigUpdated(settings config.AppSettings) {

	privKeyPath, certificatePath := generatePrivKeyPathAndCertificatePath(settings.RootPath)

	var err error
	keyPEMBlock, err := os.ReadFile(privKeyPath)
	var certPEMBlock []byte
	var cert tls.Certificate
	var privKey crypto.PrivateKey
	var hashCode []byte
	var x509Cert *x509.Certificate
	if err == nil {
		certPEMBlock, err = os.ReadFile(certificatePath)
	}

	if err == nil {
		ns.locker.RLock()
		hash := sha512.New()
		hash.Write(certPEMBlock)
		hash.Write(keyPEMBlock)
		hashCode = hash.Sum(nil)
		if ns.hashCode != nil && bytes.Equal(ns.hashCode, hashCode) {
			return
		}
		ns.locker.RUnlock()

		cert, err = tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	}

	if err == nil {
		block, _ := pem.Decode(keyPEMBlock)
		privKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	}

	if err == nil {
		x509Cert, err = x509.ParseCertificate(cert.Certificate[0])
	}

	if err != nil {
		err = ns.GenerateWithAppSettings(settings)
		if err == nil {
			return
		}
	}

	ns.locker.Lock()
	defer ns.locker.Unlock()
	v := ns.version
	defer func(v uint8) {
		if err != nil && ns.hashCode != nil {
			ns.version++
			ns.hashCode = nil
		}
		if v+1 != ns.version {
			return
		}
		ns.onUpdated()
	}(v)
	if err != nil {
		return
	}

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(x509Cert.PublicKey)
	if err == nil {
		ns.version++
		ns.nodeId = pubKeyBytes
		ns.pubKey = x509Cert.PublicKey
		ns.privKey = privKey
		ns.cert = cert
		ns.hashCode = hashCode
	}

}

func (ns *nodeSettings) onUpdated() {
	ns.registryLocker.RLock()
	registry := ns.registry
	ns.registryLocker.RUnlock()
	if registry == nil {
		return
	}
	listeners := runtime.ModulesForType[NodeSettingsListener](registry)
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

func (ns *nodeSettings) GenerateWithAppSettings(settings config.AppSettings) error {

	caPrivkey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&caPrivkey.PublicKey)
	if err != nil {
		return err
	}

	// generate certificate
	max := new(big.Int).Lsh(big.NewInt(1), 128)   //把 1 左移 128 位，返回给 big.Int
	serialNumber, _ := rand.Int(rand.Reader, max) //返回在 [0, max) 区间均匀随机分布的一个随机值
	template := &x509.Certificate{
		SerialNumber:          serialNumber, // SerialNumber 是 CA 颁布的唯一序列号，在此使用一个大随机数来代表它
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(100, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}, // 典型用法是指定叶子证书中的公钥的使用目的。它包括一系列的OID，每一个都指定一种用途。例如{id pkix 31}表示用于服务器端的TLS/SSL连接；{id pkix 34}表示密钥可以用于保护电子邮件。
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,                      // 指定了这份证书包含的公钥可以执行的密码操作，例如只能用于签名，但不能用来加密
		IsCA:                  true,                                                                       // 指示证书是不是ca证书
		BasicConstraintsValid: true,                                                                       // 指示证书是不是ca证书
	}
	certDer, err := x509.CreateCertificate(rand.Reader, template, template, &caPrivkey.PublicKey, caPrivkey)
	if err != nil {
		return err
	}

	privKeyPKCS8, err := x509.MarshalPKCS8PrivateKey(caPrivkey)
	if err != nil {
		return err
	}

	keyPEMBlock := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privKeyPKCS8})
	certPEMBlock := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDer})

	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return err
	}

	ns.locker.Lock()
	defer ns.locker.Unlock()

	privKeyPath, certificatePath := generatePrivKeyPathAndCertificatePath(settings.RootPath)
	err = os.MkdirAll(path.Dir(privKeyPath), 0750)
	if err != nil {
		return err
	}
	keyFile, err := os.Create(privKeyPath)
	if err != nil {
		return err
	}
	defer keyFile.Close()
	_, err = keyFile.Write(keyPEMBlock)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Dir(certificatePath), 0750)
	if err != nil {
		return err
	}
	certFile, err := os.Create(certificatePath)
	if err != nil {
		return err
	}
	defer certFile.Close()
	_, err = certFile.Write(certPEMBlock)
	if err != nil {
		return err
	}

	hash := sha512.New()
	hash.Write(certPEMBlock)
	hash.Write(keyPEMBlock)
	hashCode := hash.Sum(nil)

	ns.nodeId = pubKeyBytes
	ns.pubKey = &caPrivkey.PublicKey
	ns.privKey = caPrivkey
	ns.cert = cert
	ns.hashCode = hashCode

	ns.onUpdated()
	return err
}

func generatePrivKeyPathAndCertificatePath(rootPath string) (string, string) {

	return path.Join(rootPath, "key.pem"), path.Join(rootPath, "cert.pem")
}

type NodeManager interface {
	TraverseNodeID(func(NodeID) error) error
	TraverseNode(NodeID, func(Node) bool)
	Search(NodeID) []Node
	Delete(Node)
	SearchOrStore(Node) (Node, bool)
	Count(NodeID) int
}

type nodeManager struct {
	locker sync.RWMutex
	matrix [][]Node
}

func (mgr *nodeManager) TraverseNodeID(traverse func(NodeID) error) error {
	mgr.locker.RLock()
	defer mgr.locker.RUnlock()
	for _, nodeArr := range mgr.matrix {
		if len(nodeArr) > 0 {
			if err := traverse(nodeArr[0].ID()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (mgr *nodeManager) compareWithNodeID(nodeArr []Node, nodeId NodeID) int {
	return bytes.Compare(nodeArr[0].ID(), nodeId)
}

func (mgr *nodeManager) compareWithResourceID(node Node, resourceId []byte) int {
	return bytes.Compare(node.ResourceID(), resourceId)
}

func (mgr *nodeManager) TraverseNode(nodeId NodeID, traverse func(Node) bool) {

	nodeArr := mgr.Search(nodeId)
	if len(nodeArr) <= 0 {
		return
	}

	for _, node := range nodeArr {
		if !traverse(node) {
			break
		}
	}

}

func (mgr *nodeManager) Search(nodeId NodeID) []Node {
	mgr.locker.RLock()
	defer mgr.locker.RUnlock()
	idx, ok := slices.BinarySearchFunc(mgr.matrix, nodeId, mgr.compareWithNodeID)
	if ok {
		return slices.Clone(mgr.matrix[idx])
	}
	return nil
}

func (mgr *nodeManager) Delete(node Node) {
	mgr.locker.Lock()
	defer mgr.locker.Unlock()

	idx, ok := slices.BinarySearchFunc(mgr.matrix, node.ID(), mgr.compareWithNodeID)
	if !ok {
		return
	}

	idx_, ok := slices.BinarySearchFunc(mgr.matrix[idx], node.ResourceID(), mgr.compareWithResourceID)
	if !ok {
		return
	}

	if len(mgr.matrix[idx]) <= 1 {
		mgr.matrix = slices.Delete(mgr.matrix, idx, idx+1)
		return
	}
	mgr.matrix[idx] = slices.Delete(mgr.matrix[idx], idx_, idx_+1)
}

func (mgr *nodeManager) SearchOrStore(node Node) (Node, bool) {

	mgr.locker.Lock()
	defer mgr.locker.Unlock()

	idx, ok := slices.BinarySearchFunc(mgr.matrix, node.ID(), mgr.compareWithNodeID)

	if !ok {
		mgr.matrix = slices.Insert(mgr.matrix, idx, []Node{node})
		return node, false
	}

	idx_, ok := slices.BinarySearchFunc(mgr.matrix[idx], node.ResourceID(), mgr.compareWithResourceID)
	if !ok {
		mgr.matrix[idx] = slices.Insert(mgr.matrix[idx], idx_, node)
		return node, false
	}

	return mgr.matrix[idx][idx_], true

}

func (mgr *nodeManager) Count(nodeId NodeID) int {

	mgr.locker.RLock()
	defer mgr.locker.RUnlock()
	idx, ok := slices.BinarySearchFunc(mgr.matrix, nodeId, mgr.compareWithNodeID)
	if ok {
		return len(mgr.matrix[idx])
	}
	return 0
}

type NodeTripper interface {
	RoundTrip(context.Context, NodeID, io.Reader) (io.Reader, error)
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

type NodeGuard interface {
	Enabled() bool
	Access(NodeID) error
}

type NodeModule interface {
	Serve(NodeStream, Node) error
	Do(NodeID, *Request, ...NodeDoContextUpdater) (*Response, error)
	Request(NodeID, RequestName, io.Reader, ...NodeDoContextUpdater) (*Response, error)
	NodeTripper() NodeTripper
	SetNodeTripper(NodeTripper)
	NodeSettings() NodeSettings
	NodeManager() NodeManager
	ReloadModules() error
	Access(NodeID) error
	Control(Node) error
	NewResourceID(NodeType) NodeResourceID
}

type nodeModule struct {
	mgr            NodeManager
	mgrOnce        sync.Once
	settings       NodeSettings
	settingsOnce   sync.Once
	tripper        NodeTripper
	tripperLocker  sync.RWMutex
	registry       runtime.Registry
	registryLocker sync.RWMutex
	app            NodeApp
	appLocker      sync.RWMutex
	seq            uint32
	seqLocker      sync.Mutex
}

func New() interface{} {
	return &nodeModule{}
}

func (nm *nodeModule) NodeManager() NodeManager {
	nm.mgrOnce.Do(func() {
		nm.mgr = &nodeManager{}
	})
	return nm.mgr
}

func (nm *nodeModule) Init(registry runtime.Registry) error {

	nm.registryLocker.Lock()
	nm.registry = registry
	nm.registryLocker.Unlock()

	return nil
}

func (nm *nodeModule) Defer() error {
	nm.registryLocker.RLock()
	registry := nm.registry
	nm.registryLocker.RUnlock()
	if registry == nil {
		return constant.ErrUnavailable
	}

	return nm.ReloadModules()
}

func (nm *nodeModule) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[NodeAppModule](),
		reflect.TypeFor[NodeAppModuleProvider](),
		reflect.TypeFor[NodeGuard](),
	}
}

func (nm *nodeModule) Components() []bootstrap.Component {
	return []bootstrap.Component{
		bootstrap.NewComponent[NodeModule](nm, bootstrap.ComponentExternalScope),
	}
}

func (nm *nodeModule) Modules() []interface{} {
	return []interface{}{
		nm.NodeSettings(),
		nm.NodeManager(),
	}
}

func (nm *nodeModule) NodeTripper() NodeTripper {
	nm.tripperLocker.RLock()
	if nm.tripper != nil {
		defer nm.tripperLocker.RUnlock()
		return nm.tripper
	}
	nm.tripperLocker.RUnlock()
	return nm
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

	ctx := NewAppContext()
	ctx.Set(ContextNode, target)
	err := UnmarshalRequest(stream, ctx.Request())
	if err == nil {
		if app == nil {
			err = constant.ErrUnavailable
		} else {
			err = app.Run(ctx, nil)
		}
	}

	if err != nil {
		ctx.ThrowError(constant.CodeInternalError, err)
	}

	if ctx.Code() < 0 {
		ctx.ThrowError(constant.CodeNotFound, constant.ErrNotFound)
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

func (nm *nodeModule) Do(nodeId NodeID, request *Request, updaters ...NodeDoContextUpdater) (*Response, error) {

	doContext := NodeDoContext{}
	for _, updater := range updaters {
		updater(doContext)
	}
	if doContext.ctx == nil {
		doContext.ctx = context.Background()
	}

	ctx := doContext.ctx
	tripper := nm.NodeTripper()
	reqReader := MarshalRequest(request)

	resReader, err := tripper.RoundTrip(ctx, nodeId, reqReader)

	if err != nil {
		return nil, err
	}

	response := &Response{}
	InitResponse(response)
	err = UnmarshalResponse(resReader, response)

	return response, err
}

func (nm *nodeModule) Request(nodeId NodeID, name RequestName, body io.Reader, updaters ...NodeDoContextUpdater) (*Response, error) {
	request := NewRequest(name, body)
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
	app := NewApp()

	err := runtime.TraverseRegistry(registry, func(module NodeAppModule) error {
		return module.SetupToNode(app)
	})
	if err != nil {
		return err
	}

	err = runtime.TraverseRegistry(registry, func(module NodeAppModuleProvider) error {
		var setupErr error
		var router NodeRouter
		if scopeModule, ok := module.(NodeScopeModule); ok {
			scope := scopeModule.NodeScope()
			router = app.Route(scope)
		} else {
			router = app.Route(nil)
		}

		for _, m := range module.NodeAppModules() {
			setupErr = m.SetupToNode(router)
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

func (nm *nodeModule) RoundTrip(ctx context.Context, nodeId NodeID, reqReader io.Reader) (reader io.Reader, err error) {

	mgr := nm.NodeManager()
	mgr.TraverseNode(nodeId, func(node Node) bool {
		reader, err = node.Do(ctx, reqReader)
		if errors.Is(err, constant.ErrNodeClosed) {
			mgr.Delete(node)
		}
		return err != nil && !errors.Is(err, ctx.Err())
	})

	if err == nil && reader == nil {
		err = constant.ErrNodeClosed
	}
	return
}

func (nm *nodeModule) Access(nodeId NodeID) error {
	nm.registryLocker.RLock()
	defer nm.registryLocker.RUnlock()

	return runtime.TraverseRegistry(nm.registry, func(module NodeGuard) error {
		if !module.Enabled() {
			return nil
		}
		return module.Access(nodeId)
	})
}

func (nm *nodeModule) Control(node Node) error {
	err := nm.Access(node.ID())

	if err == nil {
		mgr := nm.NodeManager()
		_, ok := mgr.SearchOrStore(node)
		if ok {
			err = constant.ErrConflict
		}
	}

	return err
}

func (nm *nodeModule) NewResourceID(nodeType NodeType) NodeResourceID {

	nm.seqLocker.Lock()
	nm.seq++
	seq := nm.seq
	nm.seqLocker.Unlock()

	resourceId := make([]byte, 13)
	resourceId[0] = byte(nodeType)
	binary.BigEndian.PutUint64(resourceId[1:], uint64(time.Now().Unix()))
	binary.BigEndian.PutUint32(resourceId[9:], seq)

	return NodeResourceID(resourceId)
}
