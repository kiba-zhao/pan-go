package peer_test

import (
	"bytes"
	"crypto/rand"
	"io"
	mrand "math/rand"
	"testing"

	"treasure/peer"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestContext ...
func TestContext(t *testing.T) {
	t.Run("NewContext and Respond success", func(t *testing.T) {

		method := make([]byte, 32)
		rand.Read(method)
		body := make([]byte, 64)
		rand.Read(body)
		bodyReader := bytes.NewReader(body)
		headerName := make([]byte, 8)
		rand.Read(headerName)
		headerValue := make([]byte, 16)
		rand.Read(headerValue)

		header := peer.NewHeaderSegment(headerName, headerValue)
		req := peer.NewRequest(method, bodyReader, header)
		reader, err := peer.MarshalRequest(req)
		if err != nil {
			t.Fatal(err)
		}

		responseBodyBytes := make([]byte, 32)
		rand.Read(responseBodyBytes)
		responseBody := bytes.NewReader(responseBodyBytes)
		responseReader, writer := io.Pipe()

		stream := new(peer.RWCNodeStream)
		stream.Reader = reader
		stream.Writer = writer
		stream.Closer = writer

		id, err := uuid.NewUUID()
		if err != nil {
			t.Fatal(err)
		}
		peerId := peer.PeerId(id)
		ctx, err := peer.NewContext(stream, peerId)
		if err != nil {
			t.Fatal(err)
		}

		headers := ctx.Headers()
		bodyBytes, err := io.ReadAll(ctx.Body())
		if err != nil {
			t.Fatal(err)
		}

		go func() {
			err = ctx.Respond(responseBody)
			if err != nil {
				t.Fatal(err)
			}
		}()

		res := new(peer.Response)
		err = peer.UnmarshalResponse(responseReader, res)
		if err != nil {
			t.Fatal(err)
		}
		resBodyBytes, err := io.ReadAll(res.Body())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, peerId, ctx.PeerId(), "PeerId should be same")
		assert.Equal(t, method, ctx.Method(), "Method should be same")
		assert.Len(t, headers, 1, "Headers should has one")
		assert.Equal(t, headerName, headers[0].Name(), "Header Name should be same")
		assert.Equal(t, headerValue, headers[0].Value(), "Header Value should be same")
		assert.Equal(t, body, bodyBytes, "Body should be same")
		assert.Equal(t, 0, res.Code(), "Response code should be 0")
		assert.Equal(t, responseBodyBytes, resBodyBytes, "Response body should be same")

	})

	t.Run("Context ThrowError", func(t *testing.T) {

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

		responseReader, writer := io.Pipe()
		stream := new(peer.RWCNodeStream)
		stream.Reader = reader
		stream.Writer = writer
		stream.Closer = writer

		id, err := uuid.NewUUID()
		if err != nil {
			t.Fatal(err)
		}
		peerId := peer.PeerId(id)
		ctx, err := peer.NewContext(stream, peerId)
		if err != nil {
			t.Fatal(err)
		}

		code := mrand.Intn(1000)
		message := "Test Error"
		headerName := make([]byte, 8)
		rand.Read(headerName)
		headerValue := make([]byte, 16)
		rand.Read(headerValue)

		header := peer.NewHeaderSegment(headerName, headerValue)
		go func() {
			err = ctx.ThrowError(code, message, header)
			if err != nil {
				t.Fatal(err)
			}
		}()

		res := new(peer.Response)
		err = peer.UnmarshalResponse(responseReader, res)
		if err != nil {
			t.Fatal(err)
		}
		resBodyBytes, err := io.ReadAll(res.Body())
		if err != nil {
			t.Fatal(err)
		}
		headers := res.Headers()

		assert.Equal(t, code, res.Code(), "Response code should be same")
		assert.Equal(t, []byte(message), resBodyBytes, "Response body should be same")
		assert.Len(t, headers, 1, "Response headers should has one")
		assert.Equal(t, headerName, headers[0].Name(), "Response header name should be same")
		assert.Equal(t, headerValue, headers[0].Value(), "Response header value should be same")

	})

}
