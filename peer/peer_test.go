package peer_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"

	"sync"

	"io"
	mrand "math/rand"
	"testing"

	// "pan/core"
	coreMocked "pan/mocks/pan/core"
	mocked "pan/mocks/pan/peer"
	"pan/peer"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestPeer ...
func TestPeer(t *testing.T) {
	t.Run("Attach,Detach and Connect", func(t *testing.T) {
		t.Skip("To be invalidated")
		// addr := make([]byte, 32)
		// rand.Read(addr)

		// generator := new(mocked.MockPeerIdGenerator)
		// routeRepo := new(mocked.MockPeerRouteRepository)
		// node := new(mocked.MockNode)

		// dialer := new(mocked.MockNodeDialer)
		// dialer.On("Type").Once().Return(peer.QUICRouteType)
		// dialer.On("Connect", addr).Once().Return(node, nil)

		// baseId := uuid.New()
		// app := new(coreMocked.MockApp[peer.Context])

		// p := peer.New(baseId[:], app, generator, routeRepo)
		// p.Attach(dialer)
		// n, err := p.Connect(peer.QUICRouteType, addr)
		// if err != nil {
		// 	t.Fatal(err)
		// }

		// assert.Equal(t, node, n, "Node should be same")

		// generator.AssertExpectations(t)
		// routeRepo.AssertExpectations(t)
		// dialer.AssertExpectations(t)
		// node.AssertExpectations(t)

		// //  Detach connect error
		// terr := errors.New("Test Error")
		// dialer.On("Connect", addr).Once().Return(nil, terr)
		// _, err = p.Connect(peer.QUICRouteType, addr)
		// assert.Equal(t, terr, err, "Error should be same")

		// //  Detach test error
		// dialer.On("Type").Once().Return(peer.QUICRouteType)
		// p.Detach(dialer)
		// _, err = p.Connect(peer.QUICRouteType, addr)
		// assert.EqualError(t, err, "Not Found node dialer", "Node dialer should not found")

	})

	t.Run("Authenticate", func(t *testing.T) {

		baseId := uuid.New()
		app := new(coreMocked.MockApp[peer.Context])
		generator := new(mocked.MockPeerIdGenerator)

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

		var wg sync.WaitGroup
		wg.Add(1)

		var peerId peer.PeerId
		var authErr error
		go func() {
			defer wg.Done()
			p := peer.New(baseId[:], app, generator, 0)
			peerId, authErr = p.Authenticate(node, peer.TestOnlyAuthenticateMode)
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
		assert.Equal(t, []byte("Authenticate"), req.Method(), "Request method should be same")
		assert.Equal(t, baseId[:], reqBody, "Request base id should be same")
		assert.Equal(t, peer.TestOnlyAuthenticateMode, req.Header([]byte("Mode"))[0], "Request mode header should be same")

		app.AssertExpectations(t)
		generator.AssertExpectations(t)
		node.AssertExpectations(t)

	})

	t.Run("Authenticate with failed", func(t *testing.T) {

		baseId := uuid.New()
		app := new(coreMocked.MockApp[peer.Context])
		generator := new(mocked.MockPeerIdGenerator)

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
			p := peer.New(baseId[:], app, generator, 0)
			_, authErr = p.Authenticate(node, peer.TestOnlyAuthenticateMode)
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
		assert.Equal(t, []byte("Authenticate"), req.Method(), "Request method should be same")
		assert.Equal(t, baseId[:], reqBody, "Request base id should be same")
		assert.Equal(t, peer.TestOnlyAuthenticateMode, req.Header([]byte("Mode"))[0], "Request mode header should be same")

		app.AssertExpectations(t)
		generator.AssertExpectations(t)
		node.AssertExpectations(t)

	})

	t.Run("AcceptAuthenticate", func(t *testing.T) {

		baseId := uuid.New()
		app := new(coreMocked.MockApp[peer.Context])
		generator := new(mocked.MockPeerIdGenerator)

		reqReader, reqWriter := io.Pipe()
		resReader, resWriter := io.Pipe()
		stream := new(TestNodeStream)
		stream.Reader = reqReader
		stream.Writer = resWriter
		stream.Closer = resWriter

		ctx := context.Background()
		node := new(mocked.MockNode)
		node.On("AcceptNodeStream", ctx).Once().Return(stream, nil)
		// _, cert, err := core.GenerateKeyAndCert()
		// if err != nil {
		// 	t.Fatal(err)
		// }
		// x509Cert, err := core.ParseCertWithPem(cert)
		// if err != nil {
		// 	t.Fatal(err)
		// }
		// node.On("Certificate").Once().Return(x509Cert)

		reqBaseId := uuid.New()
		reqPeerId := peer.PeerId(uuid.New())

		generator.On("Generate", reqBaseId[:], node).Once().Return(reqPeerId, nil)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			p := peer.New(baseId[:], app, generator, 0)
			p.AcceptAuthenticate(ctx, node)
		}()

		header := peer.NewHeaderSegment([]byte("Mode"), []byte{peer.TestOnlyAuthenticateMode})
		req := peer.NewRequest([]byte("Authenticate"), bytes.NewReader(reqBaseId[:]), header)
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

		app.AssertExpectations(t)
		generator.AssertExpectations(t)
		node.AssertExpectations(t)

	})

	t.Run("Open", func(t *testing.T) {
		t.Skip("TODO: To be implement")
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

		baseId := uuid.New()
		app := new(coreMocked.MockApp[peer.Context])
		generator := new(mocked.MockPeerIdGenerator)

		var wg sync.WaitGroup
		wg.Add(1)
		var res *peer.Response
		var rresBody []byte
		go func() {
			defer wg.Done()
			p := peer.New(baseId[:], app, generator, 0)
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

		node.AssertExpectations(t)
		app.AssertExpectations(t)
		generator.AssertExpectations(t)
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

		baseId := uuid.New()
		app := new(coreMocked.MockApp[peer.Context])
		generator := new(mocked.MockPeerIdGenerator)

		p := peer.New(baseId[:], app, generator, 0)

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

		node.AssertExpectations(t)
		app.AssertExpectations(t)
		generator.AssertExpectations(t)
	})

	t.Run("AcceptServe", func(t *testing.T) {
		t.Skip("TODO: To be implement")
	})

	t.Run("Accept", func(t *testing.T) {

		generator := new(mocked.MockPeerIdGenerator)
		node := new(mocked.MockNode)
		app := new(coreMocked.MockApp[peer.Context])
		baseId := uuid.New()
		peerId := peer.PeerId(uuid.New())
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
		})

		p := peer.New(baseId[:], app, generator, 0)
		ctx := context.Background()
		p.Accept(ctx, node, peerId)

		wg.Wait()
		assert.Equal(t, ctx, acceptCtx, "Accept context.Context should be same")
		assert.Equal(t, method, appCtx.Method(), "App Context Method should be same")
		assert.Equal(t, peerId, appCtx.PeerId(), "App Context PeerId should be same")

		generator.AssertExpectations(t)
		node.AssertExpectations(t)
		app.AssertExpectations(t)
		mockStream.AssertExpectations(t)

	})
}
