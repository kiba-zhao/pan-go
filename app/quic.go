package app

import (
	"bytes"
	"cmp"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"pan/app/cache"
	"pan/runtime"
	"slices"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

type quicNode struct {
	seq  uint
	id   NodeID
	conn quic.Connection
}

func (qn *quicNode) ID() NodeID {
	return qn.id
}

func (qn *quicNode) Close() error {
	return qn.conn.CloseWithError(quic.ApplicationErrorCode(quic.NoError), "")
}

func (qn *quicNode) HashCode() uint {
	return qn.seq
}

type quicNestNodeBucket = *cache.NestBucket[NodeID, uint, *quicNode]
type quicNodeBucket = cache.Bucket[NodeID, quicNestNodeBucket]

type quicNodeConnManager struct {
	bucket quicNodeBucket
	seq    uint
	lock   sync.Mutex
}

func (mgr *quicNodeConnManager) Search(id NodeID) (*quicNode, bool) {
	bucket, ok := mgr.bucket.Search(id)
	if !ok {
		return nil, false
	}

	size := bucket.Size()
	if size <= 0 {
		return nil, false
	}

	idx, ok := bucket.Index(mgr.seq)
	if ok {
		return bucket.At(idx)
	}

	if idx > 0 {
		return bucket.At(idx - 1)
	}
	return bucket.At(size - 1)
}

func (mgr *quicNodeConnManager) Store(conn quic.Connection) (*quicNode, error) {

	state := conn.ConnectionState()
	cert := state.TLS.PeerCertificates[0]
	pubKey, err := x509.MarshalPKIXPublicKey(cert.PublicKey)
	if err != nil {
		return nil, err
	}

	bucket, ok := mgr.bucket.Search(pubKey)
	if !ok {
		bucket = &cache.NestBucket[NodeID, uint, *quicNode]{}
		bucket.Code = pubKey
		simpleBucket := cache.NewBucket[uint, *quicNode](cmp.Compare[uint])
		bucket.Bucket = cache.WrapSyncBucket(simpleBucket)
		bucket, _ = mgr.bucket.SearchOrStore(bucket)
	}

	item := &quicNode{}
	item.id = pubKey
	item.conn = conn
	for {
		mgr.lock.Lock()
		item.seq = mgr.seq
		mgr.seq++
		mgr.lock.Unlock()

		_, ok := bucket.SearchOrStore(item)
		if !ok {
			break
		}
	}

	return item, err
}

func (mgr *quicNodeConnManager) Delete(node *quicNode) {
	bucket, ok := mgr.bucket.Search(node.ID())
	if ok {
		bucket.Delete(node)
	}
}

type quicServer struct {
	quicModule *quicModule
	locker     sync.RWMutex
	ln         *quic.Listener
	address    string
}

func (qs *quicServer) Shutdown() error {
	qs.locker.RLock()
	ln := qs.ln
	qs.locker.RUnlock()
	if ln == nil {
		return ErrUnavailable
	}

	qs.locker.Lock()
	qs.ln = nil
	qs.locker.Unlock()
	return ln.Close()
}

func (qs *quicServer) ListenAndServe() error {

	if qs.quicModule == nil || qs.quicModule.NodeModule == nil {
		return ErrUnavailable
	}

	nodeSettings := qs.quicModule.NodeModule.NodeSettings()
	if !nodeSettings.Available() {
		return ErrUnavailable
	}

	certificate := nodeSettings.Certificate()
	tlsConf := &tls.Config{Certificates: []tls.Certificate{certificate}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13}
	quicConf := &quic.Config{}
	ln, err := quic.ListenAddr(qs.address, tlsConf, quicConf)
	if err != nil {
		return err
	}
	qs.locker.Lock()
	qs.ln = ln
	qs.locker.Unlock()
	defer qs.Shutdown()

	for {
		conn, err := ln.Accept(context.Background())
		if errors.Is(err, quic.ErrServerClosed) {
			break
		}
		if err != nil {
			go conn.CloseWithError(quic.ApplicationErrorCode(quic.InternalError), err.Error())
			continue
		}
		go qs.quicModule.Serve(conn)
	}

	return err
}

func (qs *quicServer) HashCode() string {
	return qs.address
}

type QuicModule interface {
	Serve(quic.Connection) error
}

type quicModule struct {
	NodeModule    NodeModule
	connMgr       *quicNodeConnManager
	connMgrLocker sync.RWMutex

	addresses []string
	locker    sync.Locker
	sigChan   chan bool
	sigOnce   sync.Once
	hasSig    bool
}

func (qm *quicModule) setSig(sig bool) {
	if qm.hasSig {
		return
	}

	qm.sigOnce.Do(func() {
		qm.sigChan = make(chan bool, 1)
	})

	qm.hasSig = true
	qm.sigChan <- sig
}

func (qm *quicModule) OnNodeSettingsUpdated(settings NodeSettings) {
	qm.locker.Lock()
	defer qm.locker.Unlock()

	qm.setSig(true)
}

func (qm *quicModule) OnConfigUpdated(settings AppSettings) {

	qm.locker.Lock()
	defer qm.locker.Unlock()

	if slices.Equal(qm.addresses, settings.NodeAddress) {
		return
	}

	qm.addresses = settings.NodeAddress
	qm.setSig(true)
}

func (qm *quicModule) Serve(conn quic.Connection) error {

	qm.connMgrLocker.RLock()
	mgr := qm.connMgr
	qm.connMgrLocker.RUnlock()

	if mgr == nil {
		return ErrUnavailable
	}

	node_, err := mgr.Store(conn)
	if err != nil {
		return conn.CloseWithError(quic.ApplicationErrorCode(quic.InternalError), err.Error())
	}
	defer mgr.Delete(node_)
	defer node_.Close()

	nodeModule := qm.NodeModule
	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			break
		}

		go nodeModule.Serve(stream, node_)
	}

	return err
}

func (qm *quicModule) Do(nodeId NodeID, reader io.Reader, ctx context.Context) (io.Reader, error) {
	qm.connMgrLocker.RLock()
	mgr := qm.connMgr
	qm.connMgrLocker.RUnlock()

	if mgr == nil {
		return nil, ErrUnavailable
	}

	node, ok := mgr.Search(nodeId)
	if !ok {
		return nil, ErrNotFound
	}

	stream, err := node.conn.OpenStream()
	defer func() {
		if err != nil && !errors.Is(err, ctx.Err()) {
			mgr.Delete(node)
		}
	}()
	if err != nil {
		return nil, err
	}

	errChan := make(chan error)
	go func() {

		_, err = io.Copy(stream, reader)
		if err == nil {
			err = stream.Close()
		}

		errChan <- err
	}()

	select {
	case err = <-errChan:
	case <-ctx.Done():
		err = ctx.Err()
	}

	return stream, err
}

func (qm *quicModule) Lookup(nodeId NodeID, ctx context.Context) error {
	qm.connMgrLocker.RLock()
	mgr := qm.connMgr
	qm.connMgrLocker.RUnlock()
	_, ok := mgr.Search(nodeId)
	if !ok {
		return ErrNotFound
	}

	return nil
}

func (qm *quicModule) Components() []runtime.Component {
	return []runtime.Component{
		runtime.NewComponent[QuicModule](qm, runtime.ComponentExternalScope),
	}
}

func (qm *quicModule) Ready() error {
	var serverMgr cache.Bucket[string, *quicServer]
	qm.sigOnce.Do(func() {
		qm.sigChan = make(chan bool, 1)
	})

	qm.connMgrLocker.Lock()
	if qm.connMgr == nil {
		nodeBucket := cache.NewBucket[NodeID, quicNestNodeBucket](bytes.Compare)
		qm.connMgr = &quicNodeConnManager{bucket: cache.WrapSyncBucket(nodeBucket)}
	}
	qm.connMgrLocker.Unlock()

	for {
		sig := <-qm.sigChan

		qm.locker.Lock()
		qm.hasSig = false
		addresses := qm.addresses
		qm.locker.Unlock()

		if serverMgr != nil {
			for _, item := range serverMgr.Items() {
				item.Shutdown()
			}
		}

		if !sig {
			break
		}

		serverMgr_ := cache.NewBucket[string, *quicServer](cmp.Compare[string])
		serverMgr = cache.WrapSyncBucket(serverMgr_)

		for _, address := range addresses {
			server := &quicServer{
				address:    address,
				quicModule: qm,
			}

			serverMgr.Store(server)
			go func(s *quicServer) {
				for {
					err := s.ListenAndServe()
					if errors.Is(err, quic.ErrServerClosed) {
						break
					}
					time.Sleep(6 * time.Second)
				}

			}(server)
		}

	}
	return nil
}
