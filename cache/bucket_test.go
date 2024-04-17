package cache_test

import (
	"bytes"
	"cmp"
	"crypto/rand"
	"pan/cache"
	"testing"

	mocked "pan/mocks/pan/cache"

	"github.com/stretchr/testify/assert"
)

func TestBucket(t *testing.T) {

	t.Run("Size", func(t *testing.T) {

		bucket := cache.NewBucket[int, *mocked.MockHashCode[int]](cmp.Compare)
		size := bucket.Size()

		assert.Equal(t, 0, size, "Size should be 0")

		item := new(mocked.MockHashCode[int])
		item.On("HashCode").Return(1)

		bucket.Store(item)

		size = bucket.Size()
		assert.Equal(t, 1, size, "Size should be 1")

		item.AssertExpectations(t)
	})

	t.Run("HashCodes", func(t *testing.T) {

		bucket := cache.NewBucket[int, *mocked.MockHashCode[int]](cmp.Compare)
		codes := bucket.HashCodes()

		assert.Nil(t, codes, "hash codes should be nil")

		item := new(mocked.MockHashCode[int])
		item.On("HashCode").Return(1)

		bucket.Store(item)

		codes = bucket.HashCodes()
		assert.Len(t, codes, 1, "hash codes should has 1")
		assert.Equal(t, 1, codes[0], "codes[0] should be 1")

		item.AssertExpectations(t)
	})

	t.Run("Items", func(t *testing.T) {

		bucket := cache.NewBucket[int, *mocked.MockHashCode[int]](cmp.Compare)
		items := bucket.Items()

		assert.Nil(t, items, "items should be nil")

		item := new(mocked.MockHashCode[int])
		item.On("HashCode").Return(1)

		bucket.Store(item)

		items = bucket.Items()
		assert.Len(t, items, 1, "items should has 1")
		assert.Equal(t, item, items[0], "items[0] should be same")

		item.AssertExpectations(t)
	})

	t.Run("At", func(t *testing.T) {

		bucket := cache.NewBucket[int, *mocked.MockHashCode[int]](cmp.Compare)
		_, ok := bucket.At(0)

		assert.False(t, ok)

		item := new(mocked.MockHashCode[int])
		defer item.AssertExpectations(t)
		item.On("HashCode").Return(1)
		bucket.Store(item)

		item_, ok := bucket.At(0)
		assert.True(t, ok)
		assert.Equal(t, item, item_)

	})

	t.Run("SearchOrStore", func(t *testing.T) {

		hash := 1

		bucket := cache.NewBucket[int, *mocked.MockHashCode[int]](cmp.Compare)
		_, ok := bucket.Search(hash)

		assert.False(t, ok)

		item := new(mocked.MockHashCode[int])
		item.On("HashCode").Return(hash)

		bitem, ok := bucket.SearchOrStore(item)

		assert.False(t, ok, "item should not found")
		assert.Equal(t, item, bitem, "item should be same")

		item1 := new(mocked.MockHashCode[int])
		item1.On("HashCode").Return(hash)

		bitem, ok = bucket.SearchOrStore(item1)

		assert.True(t, ok, "item should be found")
		assert.Equal(t, item, bitem, "item should be same")

		item.AssertExpectations(t)
		item1.AssertExpectations(t)
	})

	t.Run("Store", func(t *testing.T) {

		hash := make([]byte, 32)
		rand.Read(hash)
		item := new(mocked.MockHashCode[[]byte])
		item.On("HashCode").Times(3).Return(hash)

		bucket := cache.NewBucket[[]byte, *mocked.MockHashCode[[]byte]](bytes.Compare)
		err := bucket.Store(item)
		secondErr := bucket.Store(item)

		assert.Nil(t, err, "Error should be nil")
		assert.EqualError(t, secondErr, "Bucket Item already existsed", "Second error should be already existed error")

		item.AssertExpectations(t)

	})

	t.Run("Swap", func(t *testing.T) {
		hash := make([]byte, 32)
		rand.Read(hash)
		item := new(mocked.MockHashCode[[]byte])
		defer item.AssertExpectations(t)
		item.On("HashCode").Times(3).Return(hash)

		bucket := cache.NewBucket[[]byte, *mocked.MockHashCode[[]byte]](bytes.Compare)
		bucket.Swap(item)
		bucket.Swap(item)

	})

	t.Run("Search And Delete", func(t *testing.T) {

		bucket := cache.NewBucket[[]byte, *mocked.MockHashCode[[]byte]](bytes.Compare)

		hash := make([]byte, 32)
		rand.Read(hash)

		_, ok := bucket.Search(hash)

		assert.False(t, ok)

		item := new(mocked.MockHashCode[[]byte])
		item.On("HashCode").Times(3).Return(hash)

		err := bucket.Store(item)
		gitem, ok := bucket.Search(hash)

		assert.Nil(t, err, "Error should be nil")
		assert.True(t, ok, "Item should be found")
		assert.Equal(t, item, gitem, "Item should be same")

		oitem := new(mocked.MockHashCode[[]byte])
		oitem.On("HashCode").Times(4).Return(hash)

		bucket.Swap(oitem)
		gitem, ok = bucket.Search(hash)

		assert.True(t, ok, "Item should be found")
		assert.Equal(t, oitem, gitem, "Item should be same")
		assert.NotEqual(t, item, oitem, "Other Item should be different")

		bucket.Delete(oitem)
		_, ok = bucket.Search(hash)

		assert.False(t, ok, "Item should be not found")

		bucket.Delete(item)
		_, ok = bucket.Search(hash)

		assert.False(t, ok, "Item should be not found")

		item.AssertExpectations(t)
		oitem.AssertExpectations(t)

	})

}
