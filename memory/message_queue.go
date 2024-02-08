package memory

import "context"

type MessageQueue[T any] interface {
	Send(ctx context.Context, messages ...T) error
	Recv(ctx context.Context) ([]T, error)
}

type messageQueueSt[T any] struct {
	channel chan []T
}

// NewMessageQueue[T any] ...
func NewMessageQueue[T any](size int) MessageQueue[T] {
	mq := new(messageQueueSt[T])
	if size > 0 {
		mq.channel = make(chan []T, size)
	} else {
		mq.channel = make(chan []T)
	}
	return mq
}

// Send ...
func (mq *messageQueueSt[T]) Send(ctx context.Context, messages ...T) (err error) {

	select {
	case mq.channel <- messages:
		err = nil
	case <-ctx.Done():
		err = ctx.Err()
	}
	return
}

// Recv ...
func (mq *messageQueueSt[T]) Recv(ctx context.Context) (messages []T, err error) {
	select {
	case messages = <-mq.channel:
	case <-ctx.Done():
		err = ctx.Err()
	}

	return
}
