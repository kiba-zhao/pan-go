package app

type Next = func() error

type HandleFunc[T any] func(ctx T, next Next) error

type HandleChain[T any] []HandleFunc[T]

// Dispatch executes a chain of handle functions in sequence.
//
// It starts execution from the handle at the specified index within the chain.
// If the index is beyond the end of the chain, it calls the next function if provided.
// Each handle function is responsible for invoking the next function to continue the chain.
//
// Parameters:
// - ctx: The context of type T passed to each handle function.
// - chain: A sequence of handle functions to execute.
// - index: The current position in the chain to start execution.
// - next: An optional function to call if the end of the chain is reached.
//
// Returns an error if any handle or the next function returns an error.

func Dispatch[T any](ctx T, chain HandleChain[T], index int, next Next) error {

	if index >= len(chain) {
		if next != nil {
			return next()
		}
		return nil
	}

	handle := chain[index]
	return handle(ctx, func() error {
		return Dispatch(ctx, chain, index+1, next)
	})
}
