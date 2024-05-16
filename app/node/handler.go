package node

type Next = func() error

type HandleFunc func(ctx Context, next Next) error

type HandleChain []HandleFunc

func Dispatch(ctx Context, chain HandleChain, index int, next Next) error {

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
