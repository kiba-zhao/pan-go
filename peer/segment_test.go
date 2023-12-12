package peer_test

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"
	"treasure/peer"

	"github.com/stretchr/testify/assert"
)

// TestSegment ...
func TestSegment(t *testing.T) {

	t.Run("CreateSegmentType and ParseSegmentType", func(t *testing.T) {
		segments := []byte{peer.HeaderSegmentType, peer.BodySegmentType}
		reader := bytes.NewReader(segments)

		headerType, err := peer.ParseSegmentType(reader)
		if err != nil {
			t.Fatal(err)
		}

		bodyType, err := peer.ParseSegmentType(reader)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, peer.HeaderSegmentType, headerType, "Parse Segment Type Should be HeaderSegmentType ")
		assert.Equal(t, peer.BodySegmentType, bodyType, "Parse Segment Type Should be BodySegmentType ")

		cr := peer.CreateSegmentType(peer.HeaderSegmentType)
		headerType, err = peer.ParseSegmentType(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, peer.HeaderSegmentType, headerType, "Create Segment Type Should be HeaderSegmentType ")

		cr = peer.CreateSegmentType(peer.BodySegmentType)
		bodyType, err = peer.ParseSegmentType(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, peer.BodySegmentType, bodyType, "Create Segment Type Should be BodySegmentType ")

	})

	t.Run("CreateHeaderSegment and ParseHeaderSegment", func(t *testing.T) {
		name := make([]byte, 32)
		rand.Read(name)
		value := make([]byte, 32)
		rand.Read(value)

		header := peer.NewHeaderSegment(name, value)

		readers := peer.CreateHeaderSegment(header)
		reader := io.MultiReader(readers...)
		pHeader, err := peer.ParseHeaderSegment(reader)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, header.Name(), pHeader.Name(), "Name Should be same")
		assert.Equal(t, header.Value(), pHeader.Value(), "Value Should be same")

	})

	t.Run("CreateHeaderSegmentField and ParseHeaderSegmentField", func(t *testing.T) {
		field := make([]byte, 64)
		rand.Read(field)
		sr, dr := peer.CreateHeaderSegmentField(field)
		reader := io.MultiReader(sr, dr)
		pField, err := peer.ParseHeaderSegmentField(reader)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, field, pField, "Field Should be same")
	})

}
