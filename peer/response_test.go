package peer_test

import (
	"bytes"
	"crypto/rand"
	"io"
	mrand "math/rand"
	"testing"
	"treasure/peer"

	"github.com/stretchr/testify/assert"
)

// TestResponse ...
func TestResponse(t *testing.T) {

	t.Run("UnmarshalResponse and MarshalResponse", func(t *testing.T) {
		code := mrand.Intn(1000)
		body := make([]byte, 64)
		rand.Read(body)
		bodyReader := bytes.NewReader(body)
		headerName := make([]byte, 8)
		rand.Read(headerName)
		headerValue := make([]byte, 16)
		rand.Read(headerValue)

		header := peer.NewHeaderSegment(headerName, headerValue)
		response := peer.NewReponse(code, bodyReader, header)
		reader, err := peer.MarshalResponse(response)
		if err != nil {
			t.Fatal(err)
		}
		res := new(peer.Response)
		err = peer.UnmarshalResponse(reader, res)
		if err != nil {
			t.Fatal(err)
		}
		headers := res.Headers()
		rBody, err := io.ReadAll(res.Body())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, code, res.Code(), "Code Should be same")
		assert.Len(t, headers, 1, "Headers Should has one")
		assert.Equal(t, headerName, headers[0].Name(), "Header Name Should be same")
		assert.Equal(t, headerValue, headers[0].Value(), "Header Value Should be same")
		assert.Equal(t, body, rBody, "Body Should be same")

	})
}
