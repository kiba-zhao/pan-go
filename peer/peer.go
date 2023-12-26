package peer

import (
	"bytes"
	"cmp"
	"context"
	"errors"

	"io"
	"net"
	"pan/core"
	"pan/memory"
	"sync"

	"github.com/google/uuid"
)

const (
	QUICNodeType = uint8(iota)
	TCPNodeType
)

const (
	NormalAuthenticateMode = uint8(iota)
	TestOnlyAuthenticateMode
	OpenAuthenticateMode
)

const (
	OnlinePeerState = PeerState(iota)
	OfflinePeerState
	UnknownPeerState
)

type PeerId = uuid.UUID

type PeerState = uint

// compareUUIDBytes ...
func compareUUID(prev, next uuid.UUID) int {
	return bytes.Compare(prev[:], next[:])
}

// comparePeerId ...
func comparePeerId(prev, next PeerId) int {
	return bytes.Compare(prev[:], next[:])
}

type Peer interface {
	Stat(id PeerId) PeerState
	Connect(dialerType uint8, addr []byte) (Node, error)
	Authenticate(node Node, mode uint8) (PeerId, error)
	AcceptAuthenticate(ctx context.Context, node Node)
	Open(id PeerId) (Node, error)
	Request(node Node, body io.Reader, method []byte, headers ...*HeaderSegment) (*Response, error)
	AcceptServe(ctx context.Context, serve NodeServe)
	Accept(ctx context.Context, node Node, peerId PeerId)
}

type PeerIdGenerator interface {
	Generate(baseId []byte, node Node) (PeerId, error)
}

type PeerPassport struct {
	*memory.BucketItem[uuid.UUID]
	VerifyBaseId bool
	VerifyPeerId bool
}

type SimplePeerIdGenerator struct {
	*memory.Bucket[uuid.UUID, *PeerPassport]
	defaultDeny bool
}

// Generate ...
func (pg *SimplePeerIdGenerator) Generate(baseId []byte, node Node) (peerId PeerId, err error) {
	space, err := uuid.FromBytes(baseId)
	if err != nil {
		return
	}

	cert := node.Certificate()
	pubKey, err := core.ExtractPublicKeyFromCert(cert)
	if err != nil {
		return
	}

	id := uuid.NewSHA1(space, pubKey)

	pass := !pg.defaultDeny
	passport := pg.GetItem(id)
	if passport.VerifyPeerId == true {
		pass = !pass
	} else {
		passport = pg.GetItem(space)
		if passport.VerifyBaseId == true {
			pass = !pass
		}
	}

	if !pass {
		err = errors.New("Deny Peer Id")
		return
	}

	peerId = PeerId(id)
	return
}

// NewPeerIdGenerator ...
func NewPeerIdGenerator(defaultDeny bool) *SimplePeerIdGenerator {

	bucket := memory.NewBucket[uuid.UUID, *PeerPassport](compareUUID)

	generator := new(SimplePeerIdGenerator)
	generator.Bucket = bucket
	generator.defaultDeny = defaultDeny

	return generator
}

type peerRoute struct {
	*memory.BucketItem[[]byte]
	rw        *sync.RWMutex
	NodeType  uint8
	Addr      []byte
	FailedNum uint8
}

type peerRouteBucket = *memory.NestBucket[PeerId, []byte, *peerRoute]

type peerNodeItem struct {
	*memory.BucketItem[uint32]
	node Node
}

type peerNodeBucket = *memory.NestBucket[PeerId, uint32, *peerNodeItem]

type peerSt struct {
	dialers      *memory.Map[uint8, NodeDialer]
	handshakes   *memory.Map[uint8, NodeHandshake]
	generator    PeerIdGenerator
	router       *memory.Bucket[PeerId, peerRouteBucket]
	bucket       *memory.Bucket[PeerId, peerNodeBucket]
	baseId       uuid.UUID
	app          core.App[Context]
	maxFailedNum uint8
}

// Stat ...
func (p *peerSt) Stat(id PeerId) PeerState {
	nodeBucket := p.bucket.GetItem(id)
	if nodeBucket != nil && nodeBucket.Count() > 0 {
		return OnlinePeerState
	}
	routeBucket := p.router.GetItem(id)
	if routeBucket != nil && routeBucket.Count() > 0 {
		return UnknownPeerState
	}
	return OfflinePeerState
}

func (p *peerSt) Authenticate(node Node, mode uint8) (peerId PeerId, err error) {

	body := bytes.NewReader(p.baseId[:])
	method := append([]byte("Authenticate"), mode)
	headers := make([]*HeaderSegment, 0)
	if mode == NormalAuthenticateMode {
		p.handshakes.Range(func(key uint8, value NodeHandshake) bool {
			header := NewHeaderSegment([]byte{key}, value.Handshake())
			headers = append(headers, header)
			return true
		})
	}

	res, err := p.Request(node, body, method, headers...)
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

	peerId, err = p.generator.Generate(resBody, node)
	if err != nil {
		return
	}

	if mode != TestOnlyAuthenticateMode {
		go p.Accept(context.Background(), node, peerId)
	}

	if mode == OpenAuthenticateMode {
		return
	}

	route := new(peerRoute)
	route.Addr = node.Addr()
	route.NodeType = node.Type()
	route.FailedNum = 0
	route.rw = new(sync.RWMutex)
	code := make([]byte, 0)
	code = append(code, route.NodeType)
	code = append(code, route.Addr...)
	route.BucketItem = memory.NewBucketItem[[]byte](code)

	routeBucket := memory.NewNestBucket[PeerId, []byte, *peerRoute](peerId, bytes.Compare)
	routeBucket, _ = p.router.GetOrAddItem(routeBucket)
	routeItem, ok := routeBucket.GetOrAddItem(route)
	if ok {
		routeItem.rw.Lock()
		routeItem.FailedNum = 0
		routeItem.rw.Unlock()
	}

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

	expected := []byte("Authenticate")
	idx := len(expected)
	method := authCtx.Method()
	if len(method) < idx {
		_ = authCtx.ThrowError(UnauthorizedErrorCode, "Unauthorized")
		goto NextAcceptAuthenticate
	}

	if !bytes.Equal(expected, method[:idx]) {
		_ = authCtx.ThrowError(UnauthorizedErrorCode, "Unauthorized")
		goto NextAcceptAuthenticate
	}

	mode := method[idx]
	body, err := io.ReadAll(authCtx.Body())

	if err != nil {
		_ = authCtx.ThrowError(BadRequestErrorCode, "Bad Request")
		goto NextAcceptAuthenticate
	}

	peerId, err := p.generator.Generate(body, node)

	if err != nil {
		_ = authCtx.ThrowError(ForbiddenErrorCode, "Forbidden")
		goto NextAcceptAuthenticate
	}

	headers := authCtx.Headers()

	if headers != nil && len(headers) > 0 {

		for _, header := range headers {
			if len(header.Name()) != 1 {
				continue
			}
			nodeType := header.Name()[0]
			handshakeNode, handshakeErr := p.Connect(nodeType, header.Value())
			if handshakeErr != nil {
				continue
			}
			p.Authenticate(handshakeNode, TestOnlyAuthenticateMode)
			handshakeNode.Close()
		}
	}

	_ = authCtx.Respond(bytes.NewReader(p.baseId[:]))

	if mode != TestOnlyAuthenticateMode {
		p.Accept(ctx, node, peerId)
	}

}

// Connect ...
func (p *peerSt) Connect(nodeType uint8, addr []byte) (node Node, err error) {

	dialer, ok := p.dialers.Load(nodeType)
	if ok == false {
		err = errors.New("Not Found node dialer")
		return
	}

	node, err = dialer.Connect(addr)
	return
}

// Open ...
func (p *peerSt) Open(peerId PeerId) (node Node, err error) {

	nodeBucket := p.bucket.GetItem(peerId)
	if nodeBucket != nil {
		nodeItems := nodeBucket.GetAll()
		if nodeItems != nil && len(nodeItems) > 0 {
			node = nodeItems[0].node
			return
		}
	}

	routeBucket := p.router.GetItem(peerId)
	if routeBucket == nil {
		err = errors.New("Not Found peer node")
		return
	}

	routeItems := routeBucket.GetAll()
	if routeItems == nil || len(routeItems) < 0 {
		err = errors.New("Not Found peer node")
		return
	}

	for _, route := range routeItems {
		node, err = p.Connect(route.NodeType, route.Addr)
		if err == nil {
			authPeerId, err := p.Authenticate(node, OpenAuthenticateMode)
			if err == nil && bytes.Equal(authPeerId[:], peerId[:]) {
				route.rw.Lock()
				route.FailedNum = 0
				route.rw.Unlock()
				break
			}
		}

		route.rw.Lock()
		route.FailedNum++
		route.rw.Unlock()

		route.rw.RLock()
		if route.FailedNum > p.maxFailedNum {
			routeBucket.RemoveItem(route)
		}
		route.rw.RUnlock()

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

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer stream.Close()
		defer wg.Done()
		_, err = io.Copy(stream, req)
	}()

	var resErr error
	go func() {
		defer stream.Close()
		defer wg.Done()
		response := new(Response)
		resErr = UnmarshalResponse(stream, response)
		if resErr == nil {
			res = response
		}
	}()

	wg.Wait()
	if res != nil {
		err = nil
	} else if err == nil && resErr != nil {
		err = resErr
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

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			stream, err := node.AcceptNodeStream(ctx)
			if err != nil {
				break
			}
			go func() {
				defer stream.Close()
				c, err := NewContext(stream, peerId)
				if err == nil {
					// TODO: close with error
					_ = p.app.Run(c)
				}
			}()

		}
	}()

	nodeBucket := memory.NewNestBucket[PeerId, uint32, *peerNodeItem](peerId, cmp.Compare[uint32])
	nodeBucket, ok := p.bucket.GetOrAddItem(nodeBucket)

	nodeItem := new(peerNodeItem)
	nodeItem.node = node
	code := uint32(0)
	if ok {
		lastItem := nodeBucket.GetLastItem()
		if lastItem != nil {
			code = lastItem.HashCode() + 1
		}

	}
	nodeItem.BucketItem = memory.NewBucketItem[uint32](code)

	nodeBucket.AddItem(nodeItem)
	defer nodeBucket.RemoveItem(nodeItem)

	wg.Wait()
}

// New ...
func New(baseId uuid.UUID, app core.App[Context], generator PeerIdGenerator, maxFailedNum uint8) Peer {

	bucket := memory.NewBucket[PeerId, peerNodeBucket](comparePeerId)
	router := memory.NewBucket[PeerId, peerRouteBucket](comparePeerId)

	dialers := memory.NewMap[uint8, NodeDialer]()
	handshakes := memory.NewMap[uint8, NodeHandshake]()

	peer := new(peerSt)
	peer.baseId = baseId
	peer.app = app
	peer.generator = generator
	peer.maxFailedNum = maxFailedNum
	peer.bucket = bucket
	peer.dialers = dialers
	peer.router = router
	peer.handshakes = handshakes

	return peer
}
