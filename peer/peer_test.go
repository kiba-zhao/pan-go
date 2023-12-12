package peer_test

import (
	"crypto/rand"
	"errors"
	"testing"

	coreMocked "treasure/mocks/treasure/core"
	mocked "treasure/mocks/treasure/peer"
	"treasure/peer"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestPeer ...
func TestPeer(t *testing.T) {
	t.Run("Attach,Detach and Connect", func(t *testing.T) {

		addr := make([]byte, 32)
		rand.Read(addr)

		repo := new(mocked.MockPeerRepository)
		routeRepo := new(mocked.MockPeerRouteRepository)
		node := new(mocked.MockNode)

		dialer := new(mocked.MockNodeDialer)
		dialer.On("Type").Once().Return(peer.QUICRouteType)
		dialer.On("Connect", addr).Once().Return(node, nil)

		baseId := uuid.New()
		app := new(coreMocked.MockApp[peer.Context])

		p := peer.New(baseId, app, repo, routeRepo)
		p.Attach(dialer)
		n, err := p.Connect(peer.QUICRouteType, addr)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, node, n, "Node should be same")

		repo.AssertExpectations(t)
		routeRepo.AssertExpectations(t)
		dialer.AssertExpectations(t)
		node.AssertExpectations(t)

		//  Detach connect error
		terr := errors.New("Test Error")
		dialer.On("Connect", addr).Once().Return(nil, terr)
		_, err = p.Connect(peer.QUICRouteType, addr)
		assert.Equal(t, terr, err, "Error should be same")

		//  Detach test error
		dialer.On("Type").Once().Return(peer.QUICRouteType)
		p.Detach(dialer)
		_, err = p.Connect(peer.QUICRouteType, addr)
		assert.EqualError(t, err, "Not Found node dialer", "Node dialer should not found")

	})

	t.Run("Authenticate", func(t *testing.T) {
		t.Skip("TODO: To be implement")
	})

	t.Run("AcceptAuthenticate", func(t *testing.T) {
		t.Skip("TODO: To be implement")
	})

	t.Run("Open", func(t *testing.T) {
		t.Skip("TODO: To be implement")
	})

	t.Run("Request", func(t *testing.T) {
		t.Skip("TODO: To be implement")
	})

	t.Run("AcceptServe", func(t *testing.T) {
		t.Skip("TODO: To be implement")
	})

	t.Run("Accept", func(t *testing.T) {
		t.Skip("TODO: To be implement")
	})
}
