package broadcast_test

import (
	"crypto/rand"
	"pan/broadcast"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPacket ...
func TestPacket(t *testing.T) {

	t.Run("MarshalPacket and ParsePacket", func(t *testing.T) {

		packet := make([]byte, 32)
		rand.Read(packet)

		payload := broadcast.MarshalPacket(packet[:16], packet[16:])

		p, size, err := broadcast.ParsePacket(payload)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, len(payload), size, "Size should be mathched")
		assert.Equal(t, packet, p, "Packet should be same")

	})

	t.Run("MarshalPacket and ParsePacket with error", func(t *testing.T) {

		packet := make([]byte, 32)
		rand.Read(packet)

		payload := broadcast.MarshalPacket(packet[:15], packet[15:])

		p, size, err := broadcast.ParsePacket(payload[:1])

		assert.EqualError(t, err, "Payload Not Enough", "parsePacket Should throw not enough error")
		assert.Equal(t, size, -1, "Size should be -1")
		assert.Nil(t, p, "Method should be nil")

		p, size, err = broadcast.ParsePacket(payload[:6])

		assert.EqualError(t, err, "Payload Not Enough", "parsePacket Should throw not enough error")
		assert.Equal(t, size, len(payload), "Size should be length of payload")
		assert.Nil(t, p, "Method should be nil")

	})

}
