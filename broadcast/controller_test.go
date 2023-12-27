package broadcast_test

import (
	"bytes"
	"crypto/rand"
	"errors"
	"net"
	"strconv"

	"pan/broadcast"
	"pan/core"
	mocked "pan/mocks/pan/broadcast"
	peerMocked "pan/mocks/pan/peer"
	"pan/peer"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
)

// TestController ...
func TestController(t *testing.T) {

	t.Run("BroadcastAlive", func(t *testing.T) {
		if testing.Short() {
			t.Skip()
		}

		baseId := uuid.New()

		quicServeInfo := new(broadcast.ServeInfo)
		quicServeInfo.Port = int32(9000)
		quicServeInfo.Type = []byte{peer.QUICNodeType}

		pr := new(peerMocked.MockPeer)
		service := broadcast.NewService(baseId, pr, quicServeInfo)

		network := new(mocked.MockNet)
		payloadMatcher := mock.MatchedBy(func(p []byte) bool {
			packet, _, err := broadcast.ParsePacket(p)
			if err != nil {
				return false
			}
			method, body, err := core.ParsePacket(packet, 0)
			if err != nil || bytes.Equal([]byte("alive"), method) == false {
				return false
			}
			msg := new(broadcast.Alive)
			err = proto.Unmarshal(body, msg)
			if err == nil && len(msg.ServeInfos) == 1 && msg.ServeInfos[0].Port == quicServeInfo.Port && bytes.Equal(quicServeInfo.Type, msg.ServeInfos[0].Type) {
				return true
			}
			return false

		})
		network.On("Write", payloadMatcher).Return(nil).Times(5)

		ctrl := broadcast.NewController(service, network)
		ctrl.BroadcastAlive()

		network.AssertExpectations(t)
		pr.AssertExpectations(t)

	})

	t.Run("BroadcastAlive with dispatch error", func(t *testing.T) {

		baseId := uuid.New()

		quicServeInfo := new(broadcast.ServeInfo)
		quicServeInfo.Port = int32(9000)
		quicServeInfo.Type = []byte{peer.QUICNodeType}

		pr := new(peerMocked.MockPeer)
		service := broadcast.NewService(baseId, pr, quicServeInfo)

		network := new(mocked.MockNet)
		payloadMatcher := mock.MatchedBy(func(p []byte) bool {
			packet, _, err := broadcast.ParsePacket(p)
			if err != nil {
				return false
			}
			method, body, err := core.ParsePacket(packet, 0)
			if err != nil || bytes.Equal([]byte("alive"), method) == false {
				return false
			}
			msg := new(broadcast.Alive)
			err = proto.Unmarshal(body, msg)
			if err == nil && len(msg.ServeInfos) == 1 && msg.ServeInfos[0].Port == quicServeInfo.Port && bytes.Equal(quicServeInfo.Type, msg.ServeInfos[0].Type) {
				return true
			}
			return false

		})
		terr := errors.New("Testing Error")
		network.On("Write", payloadMatcher).Return(terr).Once()

		ctrl := broadcast.NewController(service, network)
		defer func() {
			if err := recover(); err != nil {
				assert.Equal(t, terr, err, "Error should be same")
			}
		}()
		ctrl.BroadcastAlive()

		network.AssertExpectations(t)
		pr.AssertExpectations(t)

	})

	t.Run("BroadcastDead", func(t *testing.T) {

		if testing.Short() {
			t.Skip()
		}

		baseId := uuid.New()

		quicServeInfo := new(broadcast.ServeInfo)
		quicServeInfo.Port = int32(9000)
		quicServeInfo.Type = []byte{peer.QUICNodeType}

		token := make([]byte, 32)
		rand.Read(token)

		pr := new(peerMocked.MockPeer)
		service := broadcast.NewService(baseId, pr, quicServeInfo)

		network := new(mocked.MockNet)
		payloadMatcher := mock.MatchedBy(func(p []byte) bool {
			packet, _, err := broadcast.ParsePacket(p)
			if err != nil {
				return false
			}
			method, body, err := core.ParsePacket(packet, 0)
			if err != nil || bytes.Equal([]byte("dead"), method) == false {
				return false
			}
			msg := new(broadcast.Death)
			err = proto.Unmarshal(body, msg)
			if err != nil {
				return false
			}
			return true

		})
		network.On("Write", payloadMatcher).Return(nil).Times(2)

		ctrl := broadcast.NewController(service, network)
		ctrl.BroadcastDead()

		network.AssertExpectations(t)
		pr.AssertExpectations(t)

	})

	t.Run("BroadcastDead with dispatch error", func(t *testing.T) {

		baseId := uuid.New()

		quicServeInfo := new(broadcast.ServeInfo)
		quicServeInfo.Port = int32(9000)
		quicServeInfo.Type = []byte{peer.QUICNodeType}

		pr := new(peerMocked.MockPeer)
		service := broadcast.NewService(baseId, pr, quicServeInfo)

		network := new(mocked.MockNet)
		payloadMatcher := mock.MatchedBy(func(p []byte) bool {
			packet, _, err := broadcast.ParsePacket(p)
			if err != nil {
				return false
			}
			method, body, err := core.ParsePacket(packet, 0)
			if err != nil || bytes.Equal([]byte("dead"), method) == false {
				return false
			}
			msg := new(broadcast.Death)
			err = proto.Unmarshal(body, msg)
			if err != nil {
				return false
			}
			return true

		})

		terr := errors.New("Testing Error")
		network.On("Write", payloadMatcher).Return(terr).Once()

		ctrl := broadcast.NewController(service, network)
		defer func() {
			if err := recover(); err != nil {
				assert.Equal(t, terr, err, "Error should be same")
			}
		}()
		ctrl.BroadcastDead()

		network.AssertExpectations(t)
		pr.AssertExpectations(t)

	})

	t.Run("Handle Alive", func(t *testing.T) {

		baseId := uuid.New()

		quicServeInfo := new(broadcast.ServeInfo)
		quicServeInfo.Port = int32(9000)
		quicServeInfo.Type = []byte{peer.QUICNodeType}

		pr := new(peerMocked.MockPeer)
		service := broadcast.NewService(baseId, pr, quicServeInfo)
		network := new(mocked.MockNet)
		node := new(peerMocked.MockNode)

		peerId := peer.PeerId(uuid.New())
		ip := "127.0.0.1"
		addr := []byte(net.JoinHostPort(ip, "9000"))
		method := []byte("alive")
		body, err := service.GenerateAliveMessage()
		if err != nil {
			t.Fatal(err)

		}
		msg := new(broadcast.Alive)
		err = proto.Unmarshal(body, msg)
		if err != nil {
			t.Fatal(err)
		}

		addrString := net.JoinHostPort(ip, strconv.Itoa(int(quicServeInfo.Port)))

		udpAddr, err := net.ResolveUDPAddr("udp", addrString)
		if err != nil {
			t.Fatal(err)
		}
		quicAddr := peer.MarshalQUICAddr(udpAddr)

		pr.On("Connect", uint8(quicServeInfo.Type[0]), quicAddr).Once().Return(node, nil)
		pr.On("Authenticate", node, peer.NormalAuthenticateMode).Once().Return(peerId, nil)

		ctx := new(mocked.MockContext)
		ctx.On("Method").Once().Return(method)
		ctx.On("Addr").Once().Return(addr)
		ctx.On("Body").Once().Return(body)

		ctrl := broadcast.NewController(service, network)
		ctrl.Handle(ctx, func() error {
			assert.Fail(t, "Next should not be called")
			return nil
		})

		network.AssertExpectations(t)
		ctx.AssertExpectations(t)
		pr.AssertExpectations(t)
		node.AssertExpectations(t)

	})

	t.Run("Handle Dead", func(t *testing.T) {

		baseId := uuid.New()

		quicServeInfo := new(broadcast.ServeInfo)
		quicServeInfo.Port = int32(9000)
		quicServeInfo.Type = []byte{peer.QUICNodeType}

		ip := "127.0.0.1"
		addr := []byte(net.JoinHostPort(ip, "9000"))

		pr := new(peerMocked.MockPeer)
		service := broadcast.NewService(baseId, pr, quicServeInfo)
		network := new(mocked.MockNet)

		method := []byte("dead")
		body, err := service.GenerateDeadMessage()
		if err != nil {
			t.Fatal(err)

		}
		msg := new(broadcast.Death)
		err = proto.Unmarshal(body, msg)
		if err != nil {
			t.Fatal(err)
		}

		ctx := new(mocked.MockContext)
		ctx.On("Method").Once().Return(method)
		ctx.On("Addr").Once().Return(addr)
		ctx.On("Body").Once().Return(body)

		ctrl := broadcast.NewController(service, network)
		ctrl.Handle(ctx, func() error {
			assert.Fail(t, "Next should not be called")
			return nil
		})

		network.AssertExpectations(t)
		ctx.AssertExpectations(t)
		pr.AssertExpectations(t)

	})

	t.Run("Handle Others", func(t *testing.T) {
		baseId := uuid.New()

		quicServeInfo := new(broadcast.ServeInfo)
		quicServeInfo.Port = int32(9000)
		quicServeInfo.Type = []byte{peer.QUICNodeType}

		pr := new(peerMocked.MockPeer)
		service := broadcast.NewService(baseId, pr, quicServeInfo)
		network := new(mocked.MockNet)

		method := []byte("others")

		ctx := new(mocked.MockContext)
		ctx.On("Method").Once().Return(method)

		ctrl := broadcast.NewController(service, network)

		nextCalled := false
		ctrl.Handle(ctx, func() error {
			nextCalled = true
			return nil
		})

		network.AssertExpectations(t)
		ctx.AssertExpectations(t)
		pr.AssertExpectations(t)

		assert.True(t, nextCalled, "Next Function should be called")
	})
}
