package servlet

import (
	"bytes"
	"cmp"
	"slices"
	"sync"
)

// type AppContext = *Context
// type AppHandleFunc = HandleFunc[AppContext]
// type AppHandleChain = HandleChain[AppContext]
// type AppServlet = Servlet[AppContext]

type HandleChainItem[C any, T HandleContext] struct {
	code    C
	handles HandleChain[T]
}

type Sequence = uint16

type Servlet[T HandleContext] struct {
	*Router[T]
	routes          []*HandleChainItem[RequestName, T]
	routesRW        sync.RWMutex
	defaults        HandleChain[T]
	defaults_       []*HandleChainItem[Sequence, T]
	defaultsLocker_ sync.Mutex
	rw              sync.RWMutex
	seq_            Sequence
}

func NewServlet[T HandleContext]() *Servlet[T] {
	servlet := &Servlet[T]{}
	servlet.Router = &Router[T]{}
	servlet.Router.servlet = servlet
	servlet.Router.seq = 0

	servlet.routes = make([]*HandleChainItem[RequestName, T], 0)
	servlet.defaults_ = make([]*HandleChainItem[Sequence, T], 0)

	return servlet
}

func (servlet *Servlet[T]) compareRoute(route *HandleChainItem[RequestName, T], name RequestName) int {
	return bytes.Compare(route.code, name)
}

func (servlet *Servlet[T]) route(name RequestName, handles HandleChain[T]) {
	route := &HandleChainItem[RequestName, T]{}
	route.code = name
	route.handles = append(handles, servlet.dispatchDefaults)

	servlet.routesRW.Lock()
	defer servlet.routesRW.Unlock()
	idx, ok := slices.BinarySearchFunc(servlet.routes, name, servlet.compareRoute)
	if !ok {
		servlet.routes = slices.Insert(servlet.routes, idx, route)
	}
}

func (servlet *Servlet[T]) compareDefaults_(item *HandleChainItem[Sequence, T], seq Sequence) int {
	return cmp.Compare(item.code, seq)
}

func (servlet *Servlet[T]) setDefaults(seq Sequence, handles HandleChain[T]) {
	servlet.defaultsLocker_.Lock()
	defer servlet.defaultsLocker_.Unlock()
	idx, ok := slices.BinarySearchFunc(servlet.defaults_, seq, servlet.compareDefaults_)
	if len(handles) > 0 {
		defaultItem := &HandleChainItem[Sequence, T]{}
		defaultItem.handles = handles
		defaultItem.code = seq
		if ok {
			servlet.defaults_[idx] = defaultItem
		} else {
			servlet.defaults_ = slices.Insert(servlet.defaults_, idx, defaultItem)
		}
		return
	}

	if ok {
		servlet.defaults_ = slices.Delete(servlet.defaults_, idx, idx+1)
	}
}

func (servlet *Servlet[T]) newRoute(name RequestName) *Router[T] {
	servlet.seq_++
	return NewRouter(servlet, name)
}

func (servlet *Servlet[T]) newGroup() HandleGroup[T] {
	return servlet.newRoute(nil)
}

func (servlet *Servlet[T]) Init(extreme bool) {

	servlet.defaultsLocker_.Lock()
	defer servlet.defaultsLocker_.Unlock()

	servlet.rw.Lock()
	defer servlet.rw.Unlock()

	for _, item := range servlet.defaults_ {
		servlet.defaults = slices.Concat(servlet.defaults, item.handles)
	}

	if extreme {
		servlet.Router = nil
		servlet.defaults_ = nil
	}
}

func (servlet *Servlet[T]) Run(ctx T, next Next) error {

	name := ctx.Name()
	servlet.routesRW.RLock()
	idx, ok := slices.BinarySearchFunc(servlet.routes, name, servlet.compareRoute)
	if ok {
		defer servlet.routesRW.RUnlock()
		route := servlet.routes[idx]
		return Dispatch(ctx, route.handles, 0, next)
	}
	servlet.routesRW.RUnlock()

	return servlet.dispatchDefaults(ctx, next)
}

func (servlet *Servlet[T]) dispatchDefaults(ctx T, next Next) error {

	servlet.rw.RLock()
	if len(servlet.defaults) > 0 {
		handles := slices.Clone(servlet.defaults)
		servlet.rw.RUnlock()
		return Dispatch(ctx, handles, 0, next)
	}
	servlet.rw.RUnlock()
	if next != nil {
		return next()
	}
	return nil

}
