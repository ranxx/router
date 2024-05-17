package router

import (
	"fmt"
	"reflect"
)

// RouterType router type
type RouterType int

// const
const (
	Type  RouterType = 0
	Value RouterType = 1
)

// Handler handler
type Handler interface{}

// Middleware middleware
type Middleware interface{}

// Path path
type Path interface{}

// Extra extra
type Extra interface{}

// // Register register
// type Register interface {
// 	// Register register
// 	Register(path Path, cb Handler, extra ...Extra)
// }

// // Router router
// type Router interface {
// 	// Register router register
// 	Register(path Path, cb Handler, extra ...Extra)

// 	// Router path to router
// 	Router(path Path, args ...interface{})

// 	// Range fc params is the same as the Register params
// 	Range(fc interface{})
// }

// var _ Router = &Root{}
// var _ Register = &Root{}

// []extra
// origin
// call

type metadata struct {
	extra  []Extra
	origin interface{}
	call   reflect.Value
}

// func To[T IHandler](m *metadata) *ExecMetadata[T] {
// 	return &ExecMetadata[T]{
// 		Call:   m.call,
// 		Origin: m.origin.(T),
// 		Extra:  &m.extra,
// 	}
// }

// type Param[T IHandler] interface {
// 	Call() reflect.Value
// 	Origin() T
// 	Extra() *[]Extra
// }

// func (m *metadata) Call() reflect.Value {
// 	return m.call
// }

// func (m *metadata) Origin() T {
// 	return m.origin.(T)
// }

// func To[T IHandler] () *ExecMetadata[T] {
// 	return &ExecMetadata[T]{
// 		Call:   m.call,
// 		Origin: m.origin.(T),
// 		extra:  &m.extra,
// 	}
// }

func metadataToCustomer[T IHandler](m *metadata) Customer[T] {
	return &metadataAdapter[T]{metadata: m}
}

type metadataAdapter[T IHandler] struct {
	*metadata
}

func (p *metadataAdapter[T]) Origin() T {
	return p.metadata.origin.(T)
}

func (p *metadataAdapter[T]) Call() reflect.Value {
	return p.metadata.call
}

func (p *metadataAdapter[T]) Extra() []Extra {
	return p.metadata.extra
}

type params struct {
	origin    *R
	_handlers []*metadata
	_mids     []*metadata
	_all      []reflect.Value
}

// root root
type root struct {
	rt             RouterType
	routerExec     Exec[Handler]
	middlewareExec Exec[Middleware]
	group          *G
	routers        map[Path]*params
	groups         map[*G]struct{}
}

// NewRouter new rotuer root
func NewRouter(opts ...Option) Router {
	opt := newDefaultOptions()
	for _, v := range opts {
		v(opt)
	}
	r := &root{
		rt:             opt.Type,
		routerExec:     opt.Router,
		middlewareExec: opt.Middleware,
		routers:        map[Path]*params{},
		groups:         map[*G]struct{}{},
	}
	r.group = newG(r)
	return r
}

// Use Middleware is a function, must the func len(args) == len(Handler args)
func (r *root) Use(mids ...Middleware) Router {
	r.group.Use(mids...)
	return r.tidy()
}

// NewG new group
func (r *root) NewG(mids ...Middleware) *G {
	g := newG(r, mids...)
	r.AddG(g)
	return g
}

// NewR new route
func (r *root) NewR(path Path, gs Gs, handler Handler, handlers ...Handler) *R {
	_r := newR(path, gs, handler, handlers[:]...)
	r.AddR(_r)
	return _r
}

// Add Handler is a function, must the func len(args) == len(Middleware args)
func (r *root) Add(path Path, handler Handler, handlers ...Handler) Router {
	return r.add(NewR(path, Gs{}, handler, handlers...))
}

// AddR add r
func (r *root) AddR(rts ...*R) Router {
	for _, rt := range rts {
		r.add(rt)
	}
	return r.tidy(rts...)
}

// AddG add g
func (r *root) AddG(gs ...*G) Router {
	for _, g := range gs {
		if _, ok := r.groups[g]; ok {
			continue
		}
		r.groups[g] = struct{}{}
		g.withRoot(r)
	}
	return r.tidy()
}

// Middleware Middleware is a function, must the func len(args) == len(Handler args)
func (r *root) Middleware(mids ...Middleware) Router {
	r.group.Use(mids...)
	return r.tidy()
}

// Register Handler is a function, must the func len(args) == len(Middleware args)
func (r *root) Register(path Path, handler Handler, extra ...Extra) {
	r.add(NewR(path, Gs{}, checkHandler(handler)[0]).WithExtra(extra...))
}

// Router router
func (r *root) Router(path Path, args ...interface{}) {
	p, ok := r.routers[r.valueType(path)]
	if !ok {
		return
	}
	if !r.routerMid(p, args...) {
		return
	}
	r.router(p, args...)
}

// Handlers get handlers
func (r *root) Handlers(path Path) []reflect.Value {
	cb, ok := r.routers[r.valueType(path)]
	if !ok {
		return []reflect.Value{}
	}
	return cb._all
}

// Range fc(path, []Handler, [extra...])
func (r *root) Range(fc interface{}) {
	c := reflect.ValueOf(fc)
	for _, v := range r.routers {
		ins := make([]reflect.Value, 0, 2+len(v.origin.extra))
		ins = append(ins, reflect.ValueOf(v.origin.path))
		ins = append(ins, reflect.ValueOf(v.origin.handlers))
		for _, v := range v.origin.extra {
			ins = append(ins, reflect.ValueOf(v))
		}
		c.Call(ins)
	}
}

func (r *root) routerMid(p *params, args ...interface{}) bool {
	if r.middlewareExec != nil {
		for _, v := range p._mids {
			mt := metadataToCustomer[Middleware](v)
			if !r.middlewareExec(mt, args[:]...) {
				return false
			}
		}
		return true
	}
	ins := make([]reflect.Value, 0, len(args))
	for _, v := range args {
		ins = append(ins, reflect.ValueOf(v))
	}
	for _, v := range p._mids {
		v.call.Call(ins)
	}
	return true
}

func (r *root) router(p *params, args ...interface{}) bool {
	if r.routerExec != nil {
		for _, v := range p._handlers {
			mt := metadataToCustomer[Handler](v)
			if !r.routerExec(mt, args[:]...) {
				return false
			}
		}
		return true
	}
	ins := make([]reflect.Value, 0, len(args))
	for _, v := range args {
		ins = append(ins, reflect.ValueOf(v))
	}
	for _, v := range p._handlers {
		v.call.Call(ins)
	}
	return true
}
func (r *root) valueType(v interface{}) interface{} {
	if r.rt == Value {
		return v
	}
	t := reflect.TypeOf(v)
	return t
}

func (r *root) add(rt *R) *root {
	t := r.valueType(rt.path)
	_router, ok := r.routers[t]
	if ok {
		handlers := make([]string, 0, len(_router.origin.handlers))
		for _, v := range _router.origin.handlers {
			handlers = append(handlers, reflect.TypeOf(v).String())
		}
		panic(fmt.Sprintf("Router %s registered by handlers %v extra %v", t, handlers, _router.origin.extra))
	}

	if len(rt.handlers) > 0 {
		checkHandler(rt.handlers...)
	}

	r.routers[t] = &params{origin: rt}
	return r.tidy(rt)
}

func (r *root) tidy(rts ...*R) *root {
	fc := func(rt *params) {
		mids := make([]*metadata, 0, 3)
		handlers := make([]*metadata, 0, 3)
		all := make([]reflect.Value, 0, 3)
		for i := range r.group.mids {
			mids = append(mids, &metadata{
				extra:  r.group.extra,
				origin: r.group.mids[i],
				call:   reflect.ValueOf(r.group.mids[i]),
			})
			all = append(all, reflect.ValueOf(r.group.mids[i]))
		}
		for _, g := range rt.origin.groups {
			if _, ok := r.groups[g]; !ok {
				continue
			}
			for i := range g.mids {
				mids = append(mids, &metadata{
					extra:  g.extra,
					origin: g.mids[i],
					call:   reflect.ValueOf(g.mids[i]),
				})
				all = append(all, reflect.ValueOf(g.mids[i]))
			}
		}
		for i := range rt.origin.handlers {
			handlers = append(handlers, &metadata{
				extra:  rt.origin.extra,
				origin: rt.origin.handlers[i],
				call:   reflect.ValueOf(rt.origin.handlers[i]),
			})
			all = append(all, reflect.ValueOf(rt.origin.handlers[i]))
		}
		rt._handlers = handlers
		rt._mids = mids
		rt._all = all
	}
	for _, v := range rts {
		rt, ok := r.routers[r.valueType(v.path)]
		if !ok {
			continue
		}
		fc(rt)
	}
	if len(rts) > 0 {
		return r
	}
	for _, rt := range r.routers {
		fc(rt)
	}
	return r
}

type Router interface {
	// NewG new group with root
	NewG(mids ...Middleware) *G

	// NewR new route without root
	NewR(path Path, gs Gs, handler Handler, handlers ...Handler) *R

	// Use same as Middleware, Middleware is a function, must the func len(args) == len(Handler args)
	Use(mids ...Middleware) Router

	// Add same as Register, Handler is a function, must the func len(args) == len(Middleware args)
	Add(path Path, handler Handler, handlers ...Handler) Router

	// AddR add r
	AddR(rts ...*R) Router

	// AddG add g
	AddG(gs ...*G) Router

	// Middleware same as Use, Middleware is a function, must the func len(args) == len(Handler args)
	Middleware(mids ...Middleware) Router

	// Register same as Add, Handler is a function, must the func len(args) == len(Middleware args)
	Register(path Path, handler Handler, extra ...Extra)

	// Router router
	Router(path Path, args ...interface{})

	// Handlers get handlers
	Handlers(path Path) []reflect.Value

	// Range fc(path, []Handlers, [extra...])
	Range(fc interface{})
}
