package core_test

import (
	"crypto/rand"
	"testing"
	"treasure/core"

	mocked "treasure/mocks/treasure/core"

	"github.com/stretchr/testify/mock"
)

// RunNext ...
func RunNext(args mock.Arguments) {
	next := args.Get(1).(core.Next)
	next()
}

// TestApp ...
func TestApp(t *testing.T) {

	t.Run("Use and Run", func(t *testing.T) {

		app := core.New[core.Context]()

		ctx := new(mocked.MockContext)

		handler1 := new(mocked.MockHandler[core.Context])
		handler1.On("Handle", ctx, mock.Anything).Run(RunNext).Once()
		app.Use(handler1)

		handler2 := new(mocked.MockHandler[core.Context])
		handler2.On("Handle", ctx, mock.Anything).Run(RunNext).Once()

		handler3 := new(mocked.MockHandler[core.Context])
		handler3.On("Handle", ctx, mock.Anything).Once()

		handler4 := new(mocked.MockHandler[core.Context])
		app.Use(handler2, handler3, handler4)

		app.Run(ctx)

		ctx.AssertExpectations(t)
		handler1.AssertExpectations(t)
		handler2.AssertExpectations(t)
		handler3.AssertExpectations(t)

	})

	t.Run("UseFn and Run", func(t *testing.T) {

		app := core.New[core.Context]()
		method := make([]byte, 32)
		rand.Read(method)

		ctx := new(mocked.MockContext)
		ctx.On("Method").Return(method).Times(3)

		handler1 := new(mocked.MockHandler[core.Context])
		app.UseFn(method[:1], handler1.Handle)

		handler2 := new(mocked.MockHandler[core.Context])
		handler2.On("Handle", ctx, mock.Anything).Run(RunNext).Once()

		handler3 := new(mocked.MockHandler[core.Context])
		handler3.On("Handle", ctx, mock.Anything).Once()

		handler4 := new(mocked.MockHandler[core.Context])

		app.UseFn(method, handler2.Handle, handler3.Handle, handler4.Handle)

		app.Run(ctx)

		ctx.AssertExpectations(t)
		handler1.AssertExpectations(t)
		handler2.AssertExpectations(t)
		handler3.AssertExpectations(t)
		handler4.AssertExpectations(t)

	})

}
