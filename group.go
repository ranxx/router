package router

type Gs []*G

type G struct {
	root     *root
	mids     []Middleware
	noAppend []*R
	extra    []Extra
}

func newG(r *root, mids ...Middleware) *G {
	return (&G{
		root:     r,
		mids:     make([]Middleware, 0, len(mids)),
		noAppend: []*R{},
		extra:    make([]Extra, 0, 1),
	}).Use(mids...)
}

func NewG(mids ...Middleware) *G {
	return (&G{
		mids:     make([]Middleware, 0, len(mids)),
		noAppend: []*R{},
	}).Use(mids...)
}

func (g *G) Use(mids ...Middleware) *G {
	if len(mids) > 0 {
		g.mids = append(g.mids, checkMiddleware(mids...)...)
	}
	return g
}

func (g *G) add(r *R) *G {
	if g.root == nil {
		g.noAppend = append(g.noAppend, r)
		return g
	}
	if len(g.noAppend) > 0 {
		g.root.AddR(g.noAppend...)
		g.noAppend = []*R{}
	}
	g.root.AddR(r)
	return g
}

func (g *G) Add(path Path, handler Handler, handlers ...Handler) *G {
	r := NewR(path, Gs{g}, handler, handlers[:]...)
	return g.add(r)
}

func (g *G) AddR(rts ...*R) *G {
	for _, r := range rts {
		r.WithG(g)
		g.add(r)
	}
	return g
}

func (g *G) Root() Router {
	return g.root
}

func (g *G) WithExtra(extra ...Extra) *G {
	for _, v := range extra {
		g.extra = append(g.extra, v)
	}
	return g
}

func (g *G) withRoot(r *root) *G {
	g.root = r
	if g.root == nil {
		return g
	}
	if len(g.noAppend) > 0 {
		g.root.AddR(g.noAppend...)
		g.noAppend = []*R{}
	}
	return g
}

func (g *G) tidy() *G {
	if g.root != nil {
		g.root.tidy()
	}
	return g
}

func (g *G) rtidy() bool {
	if g.root != nil {
		g.root.tidy()
		return true
	}
	return false
}
