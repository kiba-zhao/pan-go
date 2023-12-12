package peer

import (
	"bytes"
	"context"
	"crypto/x509"
	"errors"
	"io"
	"net"
	"pan/core"
	"sync"

	"github.com/google/uuid"
)

type PeerId uuid.UUID

type Peer interface {
	Connect(dialerType uint8, addr []byte) (Node, error)
	Attach(dialer NodeDialer) error
	Detach(dialer NodeDialer)
	Authenticate(node Node) (PeerId, error)
	AcceptAuthenticate(ctx context.Context, node Node)
	Open(id PeerId) (Node, error)
	Request(node Node, body io.Reader, method []byte, headers ...*HeaderSegment) (*Response, error)
	AcceptServe(ctx context.Context, serve NodeServe)
	Accept(ctx context.Context, node Node, peerId PeerId)
}

type peerItem struct {
	peerId PeerId
	node   Node
	idx    int
}

type peerBucket struct {
	items []*peerItem
	rw    *sync.RWMutex
}

func (pb *peerBucket) selectOne(peerId PeerId) (node Node) {

	pb.rw.RLock()

	for idx := len(pb.items) - 1; idx >= 0; idx-- {
		item := pb.items[idx]
		if bytes.Equal(peerId[:], item.peerId[:]) {
			node = item.node
			break
		}
	}

	pb.rw.RUnlock()
	return
}

func (pb *peerBucket) add(node Node, peerId PeerId) *peerItem {

	item := new(peerItem)
	item.peerId = peerId
	item.node = node

	pb.rw.Lock()

	item.idx = len(pb.items)
	pb.items = append(pb.items, item)

	pb.rw.Unlock()

	return item
}

// delete ...
func (pb *peerBucket) del(item *peerItem) {

	pb.rw.Lock()

	lastIdx := len(pb.items) - 1
	if lastIdx != item.idx {
		lastItem := pb.items[lastIdx]
		lastItem.idx = item.idx
		pb.items[item.idx] = lastItem
	}
	pb.items = pb.items[:lastIdx]

	pb.rw.Unlock()
}

type peerDialer struct {
	dialerMap map[uint8]NodeDialer
	rw        *sync.RWMutex
}

// Connect ...
func (pd *peerDialer) Connect(dialerType uint8, addr []byte) (node Node, err error) {
	pd.rw.RLock()
	dialer, ok := pd.dialerMap[dialerType]
	pd.rw.RUnlock()

	if ok == false {
		err = errors.New("Not Found node dialer")
		return
	}

	node, err = dialer.Connect(addr)
	return
}

// Attach ...
func (pd *peerDialer) Attach(dialer NodeDialer) (err error) {

	t := dialer.Type()

	pd.rw.Lock()
	cdialer, ok := pd.dialerMap[t]
	if ok == true && cdialer != dialer {
		err = errors.New("Duplicate node dialer")
	} else if ok == false {
		pd.dialerMap[t] = dialer
	}
	pd.rw.Unlock()

	return
}

// Detach ...
func (pd *peerDialer) Detach(dialer NodeDialer) {
	t := dialer.Type()

	pd.rw.Lock()
	cdialer, ok := pd.dialerMap[t]
	if ok == true && cdialer == dialer {
		delete(pd.dialerMap, t)
	}
	pd.rw.Unlock()
}

type peerSt struct {
	*peerDialer
	repo      PeerRepository
	routeRepo PeerRouteRepository
	bucket    *peerBucket
	baseId    uuid.UUID
	app       core.App[Context]
}

// Authenticate ...
func (p *peerSt) Authenticate(node Node) (peerId PeerId, err error) {

	body := bytes.NewReader(p.baseId[:])
	res, err := p.Request(node, body, []byte("Authenticate"))
	if err != nil {
		return
	}

	resBody, err := io.ReadAll(res.Body())
	if err != nil {
		return
	}

	if res.IsError() {
		err = NewReponseError(res.Code(), string(resBody))
		return
	}

	baseId, err := uuid.ParseBytes(resBody)
	if err != nil {
		return
	}

	cert := node.Certificate()
	pubKey, err := x509.MarshalPKIXPublicKey(cert.PublicKey)
	if err != nil {
		return
	}

	id := uuid.NewSHA1(baseId, pubKey)
	peerId = PeerId(id)

	passport, err := p.repo.FindOne(peerId)
	if err != nil {
		return
	}

	if passport == nil || passport.Enable == false {
		err = errors.New("Forbidden")
		return
	}

	go p.Accept(context.Background(), node, peerId)

	return
}

// AcceptAuthenticate ...
func (p *peerSt) AcceptAuthenticate(ctx context.Context, node Node) {

	goto AcceptAuthenticate

NextAcceptAuthenticate:

	p.AcceptAuthenticate(ctx, node)
	return

AcceptAuthenticate:

	stream, err := node.AcceptNodeStream(ctx)
	if err != nil {
		if errors.Is(err, net.ErrClosed) == false {
			node.Close()
		}
		return
	}

	authCtx, err := NewContext(stream, PeerId(p.baseId))
	if err != nil {
		_ = authCtx.ThrowError(BadRequestErrorCode, "Bad Request")
		goto NextAcceptAuthenticate
	}
	if bytes.Equal([]byte("Authenticate"), authCtx.Method()) == false {
		_ = authCtx.ThrowError(UnauthorizedErrorCode, "Unauthorized")
		goto NextAcceptAuthenticate
	}

	body, err := io.ReadAll(authCtx.Body())
	if err != nil {
		_ = authCtx.ThrowError(BadRequestErrorCode, "Bad Request")
		goto NextAcceptAuthenticate
	}
	baseId, err := uuid.ParseBytes(body)
	if err != nil {
		_ = authCtx.ThrowError(BadRequestErrorCode, "Bad Request")
		goto NextAcceptAuthenticate
	}

	cert := node.Certificate()
	pubKey, err := x509.MarshalPKIXPublicKey(cert.PublicKey)
	if err != nil {
		_ = authCtx.ThrowError(InternalErrorCode, "Internal Error")
		goto NextAcceptAuthenticate
	}

	id := uuid.NewSHA1(baseId, pubKey)
	peerId := PeerId(id)

	passport, err := p.repo.FindOne(peerId)
	if err != nil {
		_ = authCtx.ThrowError(InternalErrorCode, "Internal Error")
		goto NextAcceptAuthenticate
	}

	if passport == nil || passport.Enable == false {
		_ = authCtx.ThrowError(ForbiddenErrorCode, "Forbidden")
		goto NextAcceptAuthenticate
	}

	_ = authCtx.Respond(bytes.NewReader(p.baseId[:]))
	p.Accept(ctx, node, peerId)
}

// Open ...
func (p *peerSt) Open(peerId PeerId) (node Node, err error) {

	node = p.bucket.selectOne(peerId)
	if node == nil {
		return
	}

	routes, err := p.routeRepo.FindByPeerIdAndEnable(peerId, true)
	if err != nil {
		return
	}
	if len(routes) <= 0 {
		err = errors.New("Not Found peer node")
		return
	}

	for _, route := range routes {
		node, err = p.peerDialer.Connect(route.Type, route.Addr)
		if err != nil {
			continue
		}
		peerId, err := p.Authenticate(node)
		if err == nil && bytes.Equal(route.PeerId[:], peerId[:]) {
			break
		}
		if err == nil {
			route.Enable = false
			_ = p.routeRepo.UpdateOne(route)
		}
		_ = node.Close()
	}

	if node == nil && err == nil {
		err = errors.New("Unable to open")
	}

	return
}

// Request ...
func (p *peerSt) Request(node Node, body io.Reader, method []byte, headers ...*HeaderSegment) (res *Response, err error) {

	stream, err := node.OpenNodeStream()
	if err != nil {
		return
	}

	request := NewRequest(method, body, headers...)
	req, err := MarshalRequest(request)
	if err != nil {
		return
	}

	_, err = io.Copy(stream, req)
	if err != nil {
		return
	}
	err = stream.Close()
	if err != nil {
		return
	}

	response := new(Response)
	err = UnmarshalResponse(stream, response)
	if err != nil {
		res = response
	}

	return
}

// AcceptServe ...
func (p *peerSt) AcceptServe(ctx context.Context, serve NodeServe) {
	for {
		node, err := serve.Accept(ctx)
		if errors.Is(err, net.ErrClosed) {
			break
		}
		if err != nil {
			panic(err)
		}

		go p.AcceptAuthenticate(ctx, node)

	}
}

// Accept
func (p *peerSt) Accept(ctx context.Context, node Node, peerId PeerId) {

	if len(peerId) <= 0 {
		panic(errors.New("Missing peer id"))
	}
	item := p.bucket.add(node, peerId)
	defer p.bucket.del(item)

	for {
		stream, err := node.AcceptNodeStream(ctx)
		if err != nil {
			if errors.Is(err, net.ErrClosed) == true {
				break
			}
			continue
		}
		c, err := NewContext(stream, peerId)
		if err != nil {
			// TODO: close with error
			stream.Close()
			continue
		}
		go p.app.Run(c)
	}

}

// New ...
func New(baseId uuid.UUID, app core.App[Context], repo PeerRepository, routeRepo PeerRouteRepository) Peer {
	bucket := new(peerBucket)
	bucket.items = make([]*peerItem, 0)
	bucket.rw = new(sync.RWMutex)

	dialer := new(peerDialer)
	dialer.dialerMap = make(map[uint8]NodeDialer)
	dialer.rw = new(sync.RWMutex)

	peer := new(peerSt)
	peer.baseId = baseId
	peer.app = app
	peer.repo = repo
	peer.routeRepo = routeRepo
	peer.bucket = bucket
	peer.peerDialer = dialer

	return peer
}
