package broadcast_test

import (
	"bytes"
	"crypto/rand"
	"errors"
	"net"
	"sync"
	"testing"

	"pan/broadcast"
	"pan/core"
	mocked "pan/mocks/pan/broadcast"
	coreMocked "pan/mocks/pan/core"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// RunNext ...
func RunNext(args mock.Arguments) {
	next := args.Get(1).(core.Next)
	next()
}

// TestApp ...
func TestApp(t *testing.T) {

	t.Run("Dispatch", func(t *testing.T) {

		terr := errors.New("Testing Error")

		method := []byte("method-1")
		body := make([]byte, 32)
		rand.Read(body)
		s, m, b := core.MarshalPacket(method, body)
		payload := broadcast.MarshalPacket(s, m, b)

		mockNet := new(mocked.MockNet)
		mockNet.On("Write", payload).Return(terr).Once()

		err := broadcast.Dispatch(method, body, mockNet)

		assert.Equal(t, terr, err, "Error should be same")

		mockNet.AssertExpectations(t)

	})

	t.Run("Accept", func(t *testing.T) {

		method := []byte("method-1")
		body := make([]byte, 32)
		rand.Read(body)
		s, m, b := core.MarshalPacket(method, body)
		packet := broadcast.MarshalPacket(s, m, b)
		addr := []byte(net.JoinHostPort("127.0.0.1", "9000"))

		mockNet := new(mocked.MockNet)
		mockNet.On("Read", mock.Anything).Once().Return(packet, addr, nil)
		mockNet.On("Read", mock.Anything).Once().Return(nil, nil, net.ErrClosed)

		app := new(coreMocked.MockApp[broadcast.Context])
		wg := sync.WaitGroup{}
		wg.Add(1)
		ctxMatcher := mock.MatchedBy(func(ctx broadcast.Context) bool {
			return bytes.Equal(ctx.Method(), method) && bytes.Equal(ctx.Body(), body) && bytes.Equal(addr, ctx.Addr()) && ctx.Net() == mockNet
		})
		app.On("Run", ctxMatcher).Once().Return(nil).Run(func(args mock.Arguments) {
			wg.Done()
		})

		broadcast.Accept(app, mockNet)
		wg.Wait()

		mockNet.AssertExpectations(t)
		app.AssertExpectations(t)

	})

}
