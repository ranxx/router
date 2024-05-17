package router

type R struct {
	path     Path
	handlers []Handler
	groups   Gs
	extra    []Extra

	_g map[*G]struct{}
}

func newR(path Path, g Gs, handler Handler, handlers ...Handler) *R {
	_g := map[*G]struct{}{}
	gs := make(Gs, 0, len(g))
	for i := range g {
		if _, ok := _g[g[i]]; ok {
			continue
		}
		gs = append(gs, g[i])
		_g[g[i]] = struct{}{}
	}
	_handlers := make([]Handler, 0, 1+len(handlers))
	_handlers = append(_handlers, handler)
	_handlers = append(_handlers, handlers[:]...)

	return &R{
		path:     path,
		handlers: checkHandler(_handlers...),
		groups:   gs,
		extra:    make([]Extra, 0, 1),
		_g:       _g,
	}
}

// NewR new r
func NewR(path Path, g Gs, handler Handler, handlers ...Handler) *R {
	return newR(path, g, handler, handlers[:]...)
}

func (r *R) WithG(gs ...*G) *R {
	for i := range gs {
		if _, ok := r._g[gs[i]]; ok {
			continue
		}
		r.groups = append(r.groups, gs[i])
		r._g[gs[i]] = struct{}{}
	}
	return r.tidy()
}

func (r *R) WithExtra(extra ...Extra) *R {
	for _, v := range extra {
		r.extra = append(r.extra, v)
	}
	return r
}

func (r *R) tidy() *R {
	for _, g := range r.groups {
		if g.rtidy() {
			break
		}
	}
	return r
}

// WithG g
func WithG(gs ...*G) Gs {
	return gs
}
