package router

import (
	"fmt"
	"reflect"
)

type metadata struct {
	extra  []Extra
	origin interface{}
	call   reflect.Value
}

type params struct {
	origin    *route
	_handlers []*metadata
	_mids     []*metadata
	_all      []reflect.Value
}

// root root
type root struct {
	rt             RouterType
	routerExec     Exec
	middlewareExec Exec
	group          *group
	routers        map[Path]*params
	groups         map[*group]struct{}
}

// NewRouter new rotuer root
func NewRouter(opts ...Option) Root {
	opt := newDefaultOptions()
	for _, v := range opts {
		v(opt)
	}
	r := &root{
		rt:             opt.ty,
		routerExec:     opt.router,
		middlewareExec: opt.middleware,
		routers:        map[Path]*params{},
		groups:         map[*group]struct{}{},
	}
	r.group = newG(r)
	return r
}

// NewGroup new group, Middleware is a function type
func (r *root) NewGroup(mids ...Middleware) Group {
	g := newG(r, mids...)
	r.addG(g)
	return g
}

// NewRouter new router, handler is a function type
func (r *root) NewRouter(path Path, handler Handler, handlers ...Handler) Router {
	_r := newR(r, path, []*group{}, handler, handlers[:]...)
	r.add(_r)
	return _r
}

// Use same as Middleware func, Middleware is a function type
func (r *root) Use(mids ...Middleware) Root {
	r.group.Use(mids...)
	return r.tidy()
}

// Add same as Register func, Handler is a function type
func (r *root) Add(path Path, handler Handler, handlers ...Handler) Root {
	return r.add(newR(r, path, []*group{}, handler, handlers...))
}

// Middleware same as Use func, Middleware is a function type
func (r *root) Middleware(mids ...Middleware) Root {
	r.group.Use(mids...)
	return r.tidy()
}

// Register same as Add func, Handler is a function type
func (r *root) Register(path Path, handler Handler, extra ...Extra) Root {
	r.add(newR(r, path, []*group{}, checkHandler(handler)[0]).withExtra(extra...))
	return r
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
			ctx := newContext(v)
			if !r.middlewareExec(ctx, args[:]...) {
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
			ctx := newContext(v)
			if !r.routerExec(ctx, args[:]...) {
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

func (r *root) add(rt *route) *root {
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

func (r *root) addG(gs ...*group) *root {
	for _, _g := range gs {
		if _, ok := r.groups[_g]; ok {
			continue
		}
		r.groups[_g] = struct{}{}
		// g.withRoot(r)
	}
	return r.tidy()
}

// AddR add r
func (r *root) addR(rts ...*route) *root {
	for _, rt := range rts {
		r.add(rt)
	}
	return r.tidy(rts...)
}

func (r *root) tidy(rts ...*route) *root {
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
