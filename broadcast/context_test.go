package broadcast_test

import (
	"crypto/rand"
	"net"
	"pan/broadcast"
	mocked "pan/mocks/pan/broadcast"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestContext ...
func TestContext(t *testing.T) {

	t.Run("Simple", func(t *testing.T) {

		method := []byte("method-1")
		body := make([]byte, 32)
		rand.Read(body)
		mockNet := new(mocked.MockNet)
		addr := []byte(net.JoinHostPort("127.0.0.1", "9000"))

		ctx := broadcast.NewContext(method, body, addr, mockNet)

		assert.Equal(t, method, ctx.Method(), "Method should be same")
		assert.Equal(t, body, ctx.Body(), "Body should be same")
		assert.Equal(t, addr, ctx.Addr(), "URN should be same")
		assert.Equal(t, mockNet, ctx.Net(), "Broadcast should be same")
	})

}
