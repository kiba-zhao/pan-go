package memory_test

import (
	"cmp"
	"pan/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBucket ...
func TestBucket(t *testing.T) {
	t.Run("PutItem,RemoveItem and FindBlockItem", func(t *testing.T) {

		bucket := memory.NewBucket[string, int](cmp.Compare[int])

		item := bucket.FindBlockItem(3)
		assert.Nil(t, item, "Item should be nil before put item")

		rItem := bucket.PutItem(3, "value 3")
		item = bucket.FindBlockItem(3)
		assert.Equal(t, "value 3", item.Value(), "Item value should be same")
		assert.False(t, item.Expired(), "Item should not be expired")

		bucket.RemoveItem(rItem)
		item = bucket.FindBlockItem(3)
		assert.Nil(t, item, "Item should be nil after remove item")
		assert.True(t, rItem.Expired(), "rItem shoud be expired")

	})

	t.Run("FindBlockItems", func(t *testing.T) {

		bucket := memory.NewBucket[string, int](cmp.Compare[int])

		bucket.PutItem(3, "value 4")
		bucket.PutItem(3, "value 3")
		bucket.PutItem(3, "value 5")
		items := bucket.FindBlockItems(3)

		assert.Len(t, items, 3, "Items length should be 3")
		assert.Equal(t, "value 4", items[0].Value(), "Items[0] value should be same")
		assert.Equal(t, "value 3", items[1].Value(), "Items[1] value should be same")
		assert.Equal(t, "value 5", items[2].Value(), "Items[2] value should be same")

	})

	// t.Run("Test", func(t *testing.T) {

	// 	bucket := memory.NewBucket[string, int](cmp.Compare[int])

	// 	item4 := bucket.PutItem(4, "value 4")
	// 	item3 := bucket.PutItem(3, "value 3")
	// 	item5 := bucket.PutItem(5, "value 5")
	// 	fmt.Println("in test simple item4 block", item4.Block())
	// 	fmt.Println("in test simple item3 block", item3.Block())
	// 	fmt.Println("in test simple item5 block", item5.Block())
	// 	fmt.Println("in test simple bucket", bucket)

	// })
}
