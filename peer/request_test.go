package peer_test

import (
	"bytes"
	"crypto/rand"
	"io"
	"pan/peer"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRequest ...
func TestRequest(t *testing.T) {
	t.Run("MarshalRequest and UnmarshalRequest", func(t *testing.T) {

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
		request := peer.NewRequest(method, bodyReader, header)

		reader, err := peer.MarshalRequest(request)
		if err != nil {
			t.Fatal(err)
		}
		req := new(peer.Request)
		err = peer.UnmarshalRequest(reader, req)
		if err != nil {
			t.Fatal(err)
		}
		headers := req.Headers()
		rBody, err := io.ReadAll(req.Body())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, method, req.Method(), "Method Should be same")
		assert.Len(t, headers, 1, "Headers Should has one")
		assert.Equal(t, headerName, headers[0].Name(), "Header Name Should be same")
		assert.Equal(t, headerValue, headers[0].Value(), "Header Value Should be same")
		assert.Equal(t, body, rBody, "Body Should be same")

	})
}
