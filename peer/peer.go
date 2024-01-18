package peer

import (
	"bytes"
	"context"
	"crypto/x509"
	"errors"

	"io"
	"net"
	"sync"

	"github.com/google/uuid"
)

const (
	QUICNodeType = NodeType(iota)
	TCPNodeType
)

const (
	OnewayAuthenticateMethod = uint8(iota)
	TwowayAuthenticateMethod
	ReverseAuthenticateMethod
	OnetimeAuthenticateMethod
)

const (
	OnlinePeerState = PeerState(iota)
	OfflinePeerState
	UnknownPeerState
)

type PeerId = uuid.UUID

type PeerState = uint

// comparePeerId ...
func comparePeerId(prev, next PeerId) int {
	return bytes.Compare(prev[:], next[:])
}

type Peer interface {
	Stat(id PeerId) PeerState
	Connect(nodeType NodeType, addr []byte) (Node, error)

	OnewayAuthenticate(node Node) (PeerId, error)
	TwowayAuthenticate(node Node) (PeerId, error)
	ReverseAuthenticate(node Node) (PeerId, error)
	OnetimeAuthenticate(node Node) (PeerId, error)

	AcceptAuthenticate(ctx context.Context, node Node)
	Open(id PeerId) (Node, error)
	Request(node Node, body io.Reader, method []byte, headers ...*HeaderSegment) (*Response, error)
	AcceptServe(ctx context.Context, serve NodeServe)
	Accept(ctx context.Context, node Node, peerId PeerId)
}

type peerSt struct {
	nodeMgr  *NodeManager
	routeMgr *RouteManager
	provider Provider
}

// NewPeer ...
func NewPeer(provider Provider) Peer {

	peer := new(peerSt)
	peer.nodeMgr = NewNodeManager()
	peer.routeMgr = NewRouteManager()
	peer.provider = provider

	return peer
}

// Stat ...
func (p *peerSt) Stat(id PeerId) PeerState {
	if p.nodeMgr.Count(id) > 0 {
		return OnlinePeerState
	}
	if p.routeMgr.Count(id) > 0 {
		return UnknownPeerState
	}
	return OfflinePeerState
}

func (p *peerSt) authenticate(node Node, method []byte, headers ...*HeaderSegment) (peerId PeerId, err error) {

	baseId := getBaseId(p.provider)
	body := bytes.NewReader(baseId[:])

	res, err := p.Request(node, body, method, headers...)
	if err != nil {
		return
	}

	resBody, err := io.ReadAll(res.Body())
	if err != nil {
		return
	}

	if res.IsError() {
		err = NewResponseError(res.Code(), string(resBody))
		return
	}

	generator := p.provider.PeerIdGenerator()
	peerId, err = generator.Generate(resBody, node)

	return
}

func (p *peerSt) OnewayAuthenticate(node Node) (PeerId, error) {

	peerId, err := p.authenticate(node, []byte{OnewayAuthenticateMethod})
	if err != nil {
		return peerId, err
	}

	go p.Accept(context.Background(), node, peerId)

	_, ok := p.routeMgr.Save(peerId, node)
	if !ok {
		p.provider.PeerEvent().OnRouteAdded(peerId)
	}

	return peerId, err
}
func (p *peerSt) TwowayAuthenticate(node Node) (PeerId, error) {
	headers := make([]*HeaderSegment, 0)

	mgr := p.provider.HandshakeManager()
	mgr.Range(func(key uint8, value NodeHandshake) bool {
		handshake := value.Handshake()
		if handshake != nil {
			header := NewHeaderSegment([]byte{key}, handshake)
			headers = append(headers, header)
		}
		return true
	})

	peerId, err := p.authenticate(node, []byte{TwowayAuthenticateMethod})
	if err != nil {
		return peerId, err
	}

	go p.Accept(context.Background(), node, peerId)

	_, ok := p.routeMgr.Save(peerId, node)
	if !ok {
		p.provider.PeerEvent().OnRouteAdded(peerId)
	}

	return peerId, err
}
func (p *peerSt) ReverseAuthenticate(node Node) (PeerId, error) {
	peerId, err := p.authenticate(node, []byte{ReverseAuthenticateMethod})
	if err != nil {
		return peerId, err
	}

	_, ok := p.routeMgr.Save(peerId, node)
	if !ok {
		p.provider.PeerEvent().OnRouteAdded(peerId)
	}

	return peerId, err
}
func (p *peerSt) OnetimeAuthenticate(node Node) (PeerId, error) {
	peerId, err := p.authenticate(node, []byte{OnetimeAuthenticateMethod})
	if err != nil {
		return peerId, err
	}

	go p.Accept(context.Background(), node, peerId)
	return peerId, err
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

	authCtx, err := NewContext(stream, uuid.Nil)
	if err != nil {
		_ = authCtx.ThrowError(BadRequestErrorCode, "Bad Request")
		goto NextAcceptAuthenticate
	}

	method := authCtx.Method()
	if len(method) != 1 || method[0] > OnetimeAuthenticateMethod {
		_ = authCtx.ThrowError(NotFoundErrorCode, "Not Found")
		goto NextAcceptAuthenticate
	}

	body, err := io.ReadAll(authCtx.Body())

	if err != nil {
		_ = authCtx.ThrowError(BadRequestErrorCode, "Bad Request")
		goto NextAcceptAuthenticate
	}

	generator := p.provider.PeerIdGenerator()
	peerId, err := generator.Generate(body, node)

	if err != nil {
		_ = authCtx.ThrowError(ForbiddenErrorCode, "Forbidden")
		goto NextAcceptAuthenticate
	}

	baseId := getBaseId(p.provider)
	err = authCtx.Respond(bytes.NewReader(baseId[:]))
	if err != nil {
		if errors.Is(err, net.ErrClosed) == false {
			node.Close()
		}
		return
	}

	if method[0] == ReverseAuthenticateMethod {
		return
	}

	go p.Accept(ctx, node, peerId)

	if method[0] != TwowayAuthenticateMethod {
		return
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
			p.ReverseAuthenticate(handshakeNode)
			handshakeNode.Close()
		}
	}

}

// Connect ...
func (p *peerSt) Connect(nodeType NodeType, addr []byte) (node Node, err error) {

	mgr := p.provider.DialerManager()
	dialer, ok := mgr.Load(nodeType)
	if ok == false {
		err = errors.New("Not Found node dialer")
		return
	}

	node, err = dialer.Connect(addr)
	return
}

// Open ...
func (p *peerSt) Open(peerId PeerId) (node Node, err error) {

	node = p.nodeMgr.Get(peerId)
	if node != nil {
		return
	}

	routes := p.routeMgr.GetAll(peerId)
	if routes == nil || len(routes) <= 0 {
		err = errors.New("Not Found peer node")
		return
	}

	threshold := getRouteDiscardThreshold(p.provider)
	for _, route := range routes {
		node, err = p.Connect(route.NodeType, route.Addr)
		if err == nil {
			authPeerId, err := p.OnetimeAuthenticate(node)
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
		if route.FailedNum > threshold {
			p.routeMgr.Delete(peerId, route)
			p.provider.PeerEvent().OnRouteRemoved(peerId)
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
		app := p.provider.App()
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
					_ = app.Run(c)
				}
			}()

		}
	}()

	nodeItem, _ := p.nodeMgr.Save(peerId, node)
	defer func() {
		p.nodeMgr.Delete(peerId, nodeItem)
		p.provider.PeerEvent().OnNodeRemoved(peerId)
	}()
	p.provider.PeerEvent().OnNodeAdded(peerId)

	wg.Wait()
}

// getBaseId ...
func getBaseId(provider Provider) uuid.UUID {
	settings := provider.Settings()
	return settings.BaseId()
}

// getCert ...
func getCert(provider Provider) *x509.Certificate {
	settings := provider.Settings()
	return settings.Cert()
}

// getRouteDiscardThreshold ...
func getRouteDiscardThreshold(provider Provider) uint8 {
	settings := provider.Settings()
	return settings.RouteDiscardThreshold()
}
