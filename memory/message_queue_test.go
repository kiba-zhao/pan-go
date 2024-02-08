package memory_test

import (
	"context"
	"pan/memory"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMessageQueue(t *testing.T) {

	t.Run("Send with Blocked", func(t *testing.T) {
		mq := memory.NewMessageQueue[int](0)
		messages := []int{1, 2, 3, 5, 7, 11, 13, 17, 23}

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
		defer cancel()
		err := mq.Send(ctx, messages...)
		assert.Error(t, err)
	})

	t.Run("Recv with Blocked", func(t *testing.T) {
		mq := memory.NewMessageQueue[int](0)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
		defer cancel()
		messages, err := mq.Recv(ctx)
		assert.Error(t, err)
		assert.Nil(t, messages)
	})

	t.Run("Send and Recv", func(t *testing.T) {
		mq := memory.NewMessageQueue[int](0)
		messages := []int{1, 2, 3, 5, 7, 11, 13, 17, 23}

		var recvMessages []int
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			err := mq.Send(context.Background(), messages...)
			if err != nil {
				t.Fatal(err)
			}
		}()
		go func() {
			defer wg.Done()
			ms, err := mq.Recv(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			recvMessages = ms
		}()
		wg.Wait()

		assert.Equal(t, messages, recvMessages)

	})

	t.Run("Send and Recv With Size", func(t *testing.T) {
		size := 3
		mq := memory.NewMessageQueue[int](size)
		messages := []int{1, 2, 3, 5, 7, 11, 13, 17, 23}

		for i := 0; i <= size; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
			err := mq.Send(ctx, messages...)
			if i == size {
				if err == nil {
					t.Fatal("Error Should not be nil")
				} else {
					break
				}
			}
			if err != nil {
				t.Fatal(err)
			}

			cancel()
		}

		for i := 0; i <= size; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
			ms, err := mq.Recv(ctx)
			if i == size {
				if err == nil {
					t.Fatal("Error Should not be nil")
				} else {
					break
				}
			}
			if err != nil {
				t.Fatal(err)
			}
			assert.Equalf(t, messages, ms, "Recv %d times should be same", i)
			cancel()
		}

	})
}
