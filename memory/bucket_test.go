package memory_test

import (
	"bytes"
	"cmp"
	"crypto/rand"
	"pan/memory"
	"testing"

	mocked "pan/mocks/pan/memory"

	"github.com/stretchr/testify/assert"
)

// TestBucket ...
func TestBucket(t *testing.T) {

	t.Run("Count", func(t *testing.T) {

		bucket := memory.NewBucket[int, *mocked.MockHashCode[int]](cmp.Compare)
		count := bucket.Count()

		assert.Equal(t, 0, count, "Count should be 0")

		item := new(mocked.MockHashCode[int])
		item.On("HashCode").Return(1)

		bucket.AddItem(item)

		count = bucket.Count()
		assert.Equal(t, 1, count, "Count should be 1")

		item.AssertExpectations(t)
	})

	t.Run("GetHashCodes", func(t *testing.T) {

		bucket := memory.NewBucket[int, *mocked.MockHashCode[int]](cmp.Compare)
		codes := bucket.GetHashCodes()

		assert.Nil(t, codes, "hash codes should be nil")

		item := new(mocked.MockHashCode[int])
		item.On("HashCode").Return(1)

		bucket.AddItem(item)

		codes = bucket.GetHashCodes()
		assert.Len(t, codes, 1, "hash codes should has 1")
		assert.Equal(t, 1, codes[0], "codes[0] should be 1")

		item.AssertExpectations(t)
	})

	t.Run("GetAll", func(t *testing.T) {

		bucket := memory.NewBucket[int, *mocked.MockHashCode[int]](cmp.Compare)
		items := bucket.GetAll()

		assert.Nil(t, items, "items should be nil")

		item := new(mocked.MockHashCode[int])
		item.On("HashCode").Return(1)

		bucket.AddItem(item)

		items = bucket.GetAll()
		assert.Len(t, items, 1, "items should has 1")
		assert.Equal(t, item, items[0], "items[0] should be same")

		item.AssertExpectations(t)
	})

	t.Run("GetLastItem", func(t *testing.T) {

		bucket := memory.NewBucket[int, *mocked.MockHashCode[int]](cmp.Compare)
		lastItem := bucket.GetLastItem()

		assert.Nil(t, lastItem, "Item should be nil")

		item := new(mocked.MockHashCode[int])
		item.On("HashCode").Return(1)
		bucket.AddItem(item)

		lastItem = bucket.GetLastItem()
		assert.Equal(t, item, lastItem, "Item should be same")

		item1 := new(mocked.MockHashCode[int])
		item1.On("HashCode").Return(2)
		bucket.AddItem(item1)

		lastItem = bucket.GetLastItem()
		assert.Equal(t, item1, lastItem, "Item1 should be same")

		item.AssertExpectations(t)
		item1.AssertExpectations(t)
	})

	t.Run("GetOrAddItem", func(t *testing.T) {

		hash := 1

		bucket := memory.NewBucket[int, *mocked.MockHashCode[int]](cmp.Compare)
		bitem := bucket.GetItem(hash)

		assert.Nil(t, bitem, "item should be nil")

		item := new(mocked.MockHashCode[int])
		item.On("HashCode").Return(hash)

		bitem, ok := bucket.GetOrAddItem(item)

		assert.False(t, ok, "item should not found")
		assert.Equal(t, item, bitem, "item should be same")

		item1 := new(mocked.MockHashCode[int])
		item1.On("HashCode").Return(hash)

		bitem, ok = bucket.GetOrAddItem(item1)

		assert.True(t, ok, "item should be found")
		assert.Equal(t, item, bitem, "item should be same")

		item.AssertExpectations(t)
		item1.AssertExpectations(t)
	})

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
