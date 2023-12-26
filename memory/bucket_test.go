package memory_test

import (
	"bytes"
	"crypto/rand"
	"pan/memory"
	"testing"

	mocked "pan/mocks/pan/memory"

	"github.com/stretchr/testify/assert"
)

// TestBucket ...
func TestBucket(t *testing.T) {

	t.Run("AddItem", func(t *testing.T) {

		hash := make([]byte, 32)
		rand.Read(hash)
		item := new(mocked.MockHashCode[[]byte])
		item.On("HashCode").Times(3).Return(hash)

		bucket := memory.NewBucket[[]byte, *mocked.MockHashCode[[]byte]](bytes.Compare)
		err := bucket.AddItem(item)
		secondErr := bucket.AddItem(item)

		assert.Nil(t, err, "Error should be nil")
		assert.EqualError(t, secondErr, "Bucket Item already existsed", "Second error should be already existed error")

		item.AssertExpectations(t)

	})

	t.Run("SetItem", func(t *testing.T) {
		hash := make([]byte, 32)
		rand.Read(hash)
		item := new(mocked.MockHashCode[[]byte])
		item.On("HashCode").Times(3).Return(hash)

		bucket := memory.NewBucket[[]byte, *mocked.MockHashCode[[]byte]](bytes.Compare)
		bucket.SetItem(item)
		bucket.SetItem(item)

		item.AssertExpectations(t)

	})

	t.Run("GetItem And RemoveItem", func(t *testing.T) {

		bucket := memory.NewBucket[[]byte, *mocked.MockHashCode[[]byte]](bytes.Compare)

		hash := make([]byte, 32)
		rand.Read(hash)

		gitem := bucket.GetItem(hash)

		assert.Nil(t, gitem, "Item should be nil")

		item := new(mocked.MockHashCode[[]byte])
		item.On("HashCode").Times(3).Return(hash)

		err := bucket.AddItem(item)
		gitem = bucket.GetItem(hash)

		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, item, gitem, "Item should be same")

		oitem := new(mocked.MockHashCode[[]byte])
		oitem.On("HashCode").Times(4).Return(hash)

		bucket.SetItem(oitem)
		gitem = bucket.GetItem(hash)

		assert.Equal(t, oitem, gitem, "Item should be same")
		assert.NotEqual(t, item, oitem, "Other Item should be different")

		bucket.RemoveItem(oitem)
		gitem = bucket.GetItem(hash)

		assert.Nil(t, gitem, "Item should be nil")

		bucket.RemoveItem(item)
		gitem = bucket.GetItem(hash)

		assert.Nil(t, gitem, "Item should be nil")

		item.AssertExpectations(t)
		oitem.AssertExpectations(t)

	})
}
