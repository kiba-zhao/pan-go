package node

type Next = func() error

type HandleFunc[T any] func(ctx T, next Next) error

type HandleChain[T any] []HandleFunc[T]

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
