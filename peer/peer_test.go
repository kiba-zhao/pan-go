package peer_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"

	"net"

	"sync"

	"io"
	mrand "math/rand"
	"testing"

	coreMocked "pan/mocks/pan/core"
	mocked "pan/mocks/pan/peer"
	"pan/peer"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestPeer ...
func TestPeer(t *testing.T) {

	t.Run("ReverseAuthenticate", func(t *testing.T) {

		baseId := uuid.New()

		provider := new(mocked.MockProvider)
		settings := new(mocked.MockSettings)
		settings.On("BaseId").Once().Return(baseId)
		provider.On("Settings").Once().Return(settings)

		generator := new(mocked.MockPeerIdGenerator)
		provider.On("PeerIdGenerator").Once().Return(generator)

		event := new(mocked.MockPeerEvent)
		provider.On("PeerEvent").Once().Return(event)

		reqReader, reqWriter := io.Pipe()
		resReader, resWriter := io.Pipe()
		stream := new(TestNodeStream)
		stream.Reader = resReader
		stream.Writer = reqWriter
		stream.Closer = reqWriter

		addr := []byte("127.0.0.1:9000")
		node := new(mocked.MockNode)
		node.On("OpenNodeStream").Once().Return(stream, nil)
		node.On("Type").Once().Return(peer.QUICNodeType)
		node.On("Addr").Once().Return(addr)

		resBaseId := uuid.New()
		resPeerId := peer.PeerId(uuid.New())

		generator.On("Generate", resBaseId[:], node).Once().Return(resPeerId, nil)

		event.On("OnRouteAdded", resPeerId).Once()

		var wg sync.WaitGroup
		wg.Add(1)

		var peerId peer.PeerId
		var authErr error
		go func() {
			defer wg.Done()
			p := peer.NewPeer(provider)
			peerId, authErr = p.ReverseAuthenticate(node)
		}()

		req := new(peer.Request)
		err := peer.UnmarshalRequest(reqReader, req)
		if err != nil {
			t.Fatal(err)
		}
		reqBody, err := io.ReadAll(req.Body())
		if err != nil {
			t.Fatal(err)
		}

		res := peer.NewReponse(0, bytes.NewReader(resBaseId[:]))
		reader, err := peer.MarshalResponse(res)
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.Copy(resWriter, reader)
		if err != nil {
			t.Fatal(err)
		}
		resWriter.Close()

		wg.Wait()

		assert.Nil(t, authErr, "Authenticate should without error")
		assert.Equal(t, resPeerId, peerId, "Peer Id should be same")
		assert.Equal(t, []byte{peer.ReverseAuthenticateMethod}, req.Method(), "Request method should be same")
		assert.Equal(t, baseId[:], reqBody, "Request base id should be same")

		provider.AssertExpectations(t)
		settings.AssertExpectations(t)
		generator.AssertExpectations(t)
		event.AssertExpectations(t)
		node.AssertExpectations(t)

	})

	t.Run("ReverseAuthenticate with failed", func(t *testing.T) {

		baseId := uuid.New()
		provider := new(mocked.MockProvider)

		settings := new(mocked.MockSettings)
		settings.On("BaseId").Once().Return(baseId)

		generator := new(mocked.MockPeerIdGenerator)
		provider.On("Settings").Once().Return(settings)
		provider.On("PeerIdGenerator").Once().Return(generator)

		reqReader, reqWriter := io.Pipe()
		resReader, resWriter := io.Pipe()
		stream := new(TestNodeStream)
		stream.Reader = resReader
		stream.Writer = reqWriter
		stream.Closer = reqWriter

		node := new(mocked.MockNode)
		node.On("OpenNodeStream").Once().Return(stream, nil)

		resBaseId := uuid.New()
		terr := errors.New("Test Error")

		generator.On("Generate", resBaseId[:], node).Once().Return(nil, terr)

		var wg sync.WaitGroup
		wg.Add(1)

		var authErr error
		go func() {
			defer wg.Done()
			p := peer.NewPeer(provider)
			_, authErr = p.ReverseAuthenticate(node)
		}()

		req := new(peer.Request)
		err := peer.UnmarshalRequest(reqReader, req)
		if err != nil {
			t.Fatal(err)
		}
		reqBody, err := io.ReadAll(req.Body())
		if err != nil {
			t.Fatal(err)
		}

		res := peer.NewReponse(0, bytes.NewReader(resBaseId[:]))
		reader, err := peer.MarshalResponse(res)
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.Copy(resWriter, reader)
		if err != nil {
			t.Fatal(err)
		}
		resWriter.Close()

		wg.Wait()

		assert.ErrorIs(t, terr, authErr, "Authenticate should be error")
		assert.Equal(t, []byte{peer.ReverseAuthenticateMethod}, req.Method())
		assert.Equal(t, baseId[:], reqBody, "Request base id should be same")

		provider.AssertExpectations(t)
		settings.AssertExpectations(t)
		generator.AssertExpectations(t)
		node.AssertExpectations(t)

	})

	t.Run("AcceptAuthenticate", func(t *testing.T) {

		baseId := uuid.New()
		provider := new(mocked.MockProvider)

		settings := new(mocked.MockSettings)
		settings.On("BaseId").Once().Return(baseId)

		generator := new(mocked.MockPeerIdGenerator)
		provider.On("Settings").Once().Return(settings)
		provider.On("PeerIdGenerator").Once().Return(generator)

		reqReader, reqWriter := io.Pipe()
		resReader, resWriter := io.Pipe()
		stream := new(TestNodeStream)
		stream.Reader = reqReader
		stream.Writer = resWriter
		stream.Closer = resWriter

		ctx := context.Background()
		node := new(mocked.MockNode)
		node.On("AcceptNodeStream", ctx).Once().Return(stream, nil)

		reqBaseId := uuid.New()
		reqPeerId := peer.PeerId(uuid.New())

		generator.On("Generate", reqBaseId[:], node).Once().Return(reqPeerId, nil)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			p := peer.NewPeer(provider)
			p.AcceptAuthenticate(ctx, node)
		}()

		req := peer.NewRequest([]byte{peer.ReverseAuthenticateMethod}, bytes.NewReader(reqBaseId[:]))
		reader, err := peer.MarshalRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.Copy(reqWriter, reader)
		if err != nil {
			t.Fatal(err)
		}
		err = reqWriter.Close()
		if err != nil {
			t.Fatal(err)
		}
		res := new(peer.Response)
		err = peer.UnmarshalResponse(resReader, res)
		if err != nil {
			t.Fatal(err)
		}
		resBody, err := io.ReadAll(res.Body())
		if err != nil {
			t.Fatal(err)
		}

		wg.Wait()

		assert.Equal(t, 0, res.Code(), "Response code should be same")
		assert.Equal(t, baseId[:], resBody, "Response base id should be same")

		provider.AssertExpectations(t)
		settings.AssertExpectations(t)
		generator.AssertExpectations(t)
		node.AssertExpectations(t)

	})

	t.Run("Open", func(t *testing.T) {

		baseId := uuid.New()
		provider := new(mocked.MockProvider)

		p := peer.NewPeer(provider)
		node, err := p.Open(peer.PeerId(baseId))

		assert.Nil(t, node, "Node should be nil")
		assert.EqualError(t, err, "Not Found peer node", "Error should be not found")

		provider.AssertExpectations(t)
	})

	t.Run("Request", func(t *testing.T) {

		method := make([]byte, 32)
		rand.Read(method)
		body := make([]byte, 64)
		rand.Read(body)
		bodyReader := bytes.NewReader(body)

		reqReader, reqWriter := io.Pipe()

		stream := new(TestNodeStream)
		stream.Closer = reqWriter
		stream.Writer = reqWriter

		resCode := mrand.Intn(1000)
		resBody := make([]byte, 64)
		rand.Read(resBody)
		resBodyReader := bytes.NewReader(resBody)

		response := peer.NewReponse(resCode, resBodyReader)
		resReader, err := peer.MarshalResponse(response)
		if err != nil {
			t.Fatal(err)
		}

		reader, writer := io.Pipe()
		stream.Reader = reader

		node := new(mocked.MockNode)
		node.On("OpenNodeStream").Once().Return(stream, nil)

		provider := new(mocked.MockProvider)

		var wg sync.WaitGroup
		wg.Add(1)
		var res *peer.Response
		var rresBody []byte
		go func() {
			defer wg.Done()
			p := peer.NewPeer(provider)
			response, err := p.Request(node, bodyReader, method)
			if err != nil {
				t.Fatal(t)
			}
			res = response
			rresBody, err = io.ReadAll(res.Body())
			if err != nil {
				t.Fatal(t)
			}

		}()

		req := new(peer.Request)
		err = peer.UnmarshalRequest(reqReader, req)
		if err != nil {
			t.Fatal(err)
		}
		reqBody, err := io.ReadAll(req.Body())
		if err != nil {
			t.Fatal(t)
		}

		_, err = io.Copy(writer, resReader)
		if err != nil {
			t.Fatal(t)
		}
		writer.Close()
		wg.Wait()

		assert.Equal(t, method, req.Method(), "Request method should be same")
		assert.Equal(t, body, reqBody, "Request Body should be same")
		assert.Equal(t, resCode, res.Code(), "Response code should be same")
		assert.Equal(t, resBody, rresBody, "Response Body should be same")

		provider.AssertExpectations(t)
		node.AssertExpectations(t)

	})

	t.Run("Request with not all read before response", func(t *testing.T) {

		method := make([]byte, 32)
		rand.Read(method)
		body := make([]byte, 64)
		rand.Read(body)
		bodyReader := bytes.NewReader(body)

		_, reqWriter := io.Pipe()

		stream := new(TestNodeStream)
		stream.Closer = reqWriter
		stream.Writer = reqWriter

		resCode := mrand.Intn(1000)
		resBody := make([]byte, 64)
		rand.Read(resBody)
		resBodyReader := bytes.NewReader(resBody)

		response := peer.NewReponse(resCode, resBodyReader)
		resReader, err := peer.MarshalResponse(response)
		if err != nil {
			t.Fatal(err)
		}
		stream.Reader = resReader

		node := new(mocked.MockNode)
		node.On("OpenNodeStream").Once().Return(stream, nil)

		provider := new(mocked.MockProvider)
		p := peer.NewPeer(provider)

		var wg sync.WaitGroup
		wg.Add(1)
		var res peer.Response
		go func() {
			defer wg.Done()
			response, err := p.Request(node, bodyReader, method)
			if err != nil {
				t.Fatal(t)
			}
			res = *response
		}()

		wg.Wait()

		rresBody, err := io.ReadAll(res.Body())
		if err != nil {
			t.Fatal(t)
		}

		wg.Wait()

		assert.Equal(t, resCode, res.Code(), "Response code should be same")
		assert.Equal(t, resBody, rresBody, "Response Body should be same")

		provider.AssertExpectations(t)
		node.AssertExpectations(t)
	})

	t.Run("AcceptServe", func(t *testing.T) {

		var wg sync.WaitGroup
		wg.Add(1)

		ctx := context.Background()
		node := new(mocked.MockNode)
		node.On("AcceptNodeStream", ctx).Once().Return(nil, net.ErrClosed).Run(func(args mock.Arguments) { defer wg.Done() })

		provider := new(mocked.MockProvider)
		serve := new(mocked.MockNodeServe)

		serve.On("Accept", ctx).Once().Return(node, nil)
		serve.On("Accept", ctx).Once().Return(nil, net.ErrClosed)

		p := peer.NewPeer(provider)
		p.AcceptServe(ctx, serve)

		wg.Wait()

		node.AssertExpectations(t)
		provider.AssertExpectations(t)
		serve.AssertExpectations(t)
	})

	t.Run("Accept", func(t *testing.T) {

		provider := new(mocked.MockProvider)
		app := new(coreMocked.MockApp[peer.Context])
		provider.On("App").Once().Return(app)

		node := new(mocked.MockNode)
		peerId := peer.PeerId(uuid.New())

		event := new(mocked.MockPeerEvent)
		provider.On("PeerEvent").Twice().Return(event)

		event.On("OnNodeAdded", peerId).Once()
		event.On("OnNodeRemoved", peerId).Once()

		stream := new(TestNodeStream)
		mockStream := new(mocked.MockNodeStream)
		stream.Closer = mockStream

		method := make([]byte, 32)
		rand.Read(method)
		body := make([]byte, 64)
		rand.Read(body)
		bodyReader := bytes.NewReader(body)

		req := peer.NewRequest(method, bodyReader)
		reader, err := peer.MarshalRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		stream.Reader = reader

		mockStream.On("Close").Once().Return(nil)
		var acceptCtx context.Context
		node.On("AcceptNodeStream", mock.Anything).Once().Return(stream, nil).Run(func(args mock.Arguments) {
			acceptCtx = args.Get(0).(context.Context)
		})
		node.On("AcceptNodeStream", mock.Anything).Once().Return(nil, errors.New("Test Error"))

		var wg sync.WaitGroup
		wg.Add(1)
		var appCtx peer.Context
		app.On("Run", mock.Anything).Once().Run(func(args mock.Arguments) {
			defer wg.Done()
			appCtx = args.Get(0).(peer.Context)
		}).Return(nil)

		p := peer.NewPeer(provider)
		ctx := context.Background()
		p.Accept(ctx, node, peerId)

		wg.Wait()
		assert.Equal(t, ctx, acceptCtx, "Accept context.Context should be same")
		assert.Equal(t, method, appCtx.Method(), "App Context Method should be same")
		assert.Equal(t, peerId, appCtx.PeerId(), "App Context PeerId should be same")

		provider.AssertExpectations(t)
		event.AssertExpectations(t)
		node.AssertExpectations(t)
		app.AssertExpectations(t)
		mockStream.AssertExpectations(t)

	})
}
