package core

type Event[T interface{}] interface {
	Attach(handler T)
	Dettach(handler T)
}
