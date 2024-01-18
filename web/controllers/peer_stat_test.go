package controllers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"pan/models"
	"pan/peer"
	"pan/services"
	"pan/web"
	"pan/web/controllers"

	"testing"

	peerMocked "pan/mocks/pan/peer"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestPeerStat ...
func TestPeerStat(t *testing.T) {

	setup := func(app *web.App, peer peer.Peer) {
		ctrl := new(controllers.PeerStatController)
		ctrl.PeerStat = new(services.PeerStatService)
		ctrl.PeerStat.Peer = peer

		ctrl.Init(app)
	}

	t.Run("GET /base/peer-stat/:id", func(t *testing.T) {

		p := new(peerMocked.MockPeer)
		app := web.NewApp()
		setup(app, p)

		id := uuid.New()
		url := fmt.Sprintf("/base/peer-stat/%s", id.String())
		p.On("Stat", id).Once().Return(peer.OnlinePeerState)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		app.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		stat := new(models.PeerStat)
		err := json.Unmarshal(w.Body.Bytes(), stat)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, id.String(), stat.ID)
		assert.Equal(t, peer.OnlinePeerState, stat.Stat)

		p.AssertExpectations(t)

	})

	t.Run("GET /base/peer-stat/:id With Bad Request", func(t *testing.T) {

		p := new(peerMocked.MockPeer)
		app := web.NewApp()
		setup(app, p)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/base/peer-stat/123", nil)
		app.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)

		p.AssertExpectations(t)

	})

}
