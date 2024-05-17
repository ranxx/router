package router

var _ Router = &route{}

type route struct {
	*root
	path     Path
	handlers []Handler
	groups   []*group
	extra    []Extra

	_g map[*group]struct{}
}

func newR(root *root, path Path, gs []*group, handler Handler, handlers ...Handler) *route {
	_g := map[*group]struct{}{}
	_gs := make([]*group, 0, len(gs))
	for i := range gs {
		if _, ok := _g[gs[i]]; ok {
			continue
		}
		gs = append(gs, gs[i])
		_g[gs[i]] = struct{}{}
	}
	_handlers := make([]Handler, 0, 1+len(handlers))
	_handlers = append(_handlers, handler)
	_handlers = append(_handlers, handlers[:]...)

	return &route{
		root:     root,
		path:     path,
		handlers: checkHandler(_handlers...),
		groups:   _gs,
		extra:    make([]Extra, 0, 1),
		_g:       _g,
	}
}

// WithExtra with extra
func (r *route) WithExtra(extra ...Extra) Router {
	return r.withExtra(extra...)
}

func (r *route) withExtra(extra ...Extra) *route {
	for _, v := range extra {
		r.extra = append(r.extra, v)
	}
	return r.tidy()
}

func (r *route) tidy() *route {
	r.root.tidy(r)
	return r
}
