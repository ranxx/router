package router

var _ Group = &group{}

type group struct {
	*root
	mids     []Middleware
	noAppend []*route
	extra    []Extra
}

func newG(root *root, mids ...Middleware) *group {
	if len(mids) > 0 {
		checkMiddleware(mids...)
	}
	return &group{
		root:     root,
		mids:     mids,
		noAppend: []*route{},
		extra:    make([]Extra, 0, 1),
	}
}

// NewRouter new Router, Handler is a function type
func (g *group) NewRouter(path Path, handler Handler, handlers ...Handler) Router {
	r := newR(g.root, path, []*group{g}, handler, handlers[:]...)
	g.tidy()
	return r
}

// Use same as Middleware func, Middleware is a function type
func (g *group) Use(mids ...Middleware) Group {
	if len(mids) > 0 {
		g.mids = append(g.mids, checkMiddleware(mids...)...)
	}
	return g.tidy()
}

// Add same as Register func, Handler is a function type
func (g *group) Add(path Path, handler Handler, handlers ...Handler) Group {
	r := newR(g.root, path, []*group{g}, handler, handlers[:]...)
	return g.add(r)
}

// Middleware same as Use func, Middleware is a function type
func (g *group) Middleware(mids ...Middleware) Group {
	return g.Use(mids...)
}

// Register same as Add func, Handler is a function type
func (g *group) Register(path Path, handler Handler, extra ...Extra) Group {
	g.add(newR(g.root, path, []*group{g}, checkHandler(handler)[0]).withExtra(extra...))
	return g.tidy()
}

// Root return Root interface
func (g *group) Root() Root {
	return g.root
}

// WithExtra with extra
func (g *group) WithExtra(extra ...Extra) Group {
	for _, v := range extra {
		g.extra = append(g.extra, v)
	}
	return g.tidy()
}

func (g *group) withRoot(root *root) *group {
	g.root = root
	if g.root == nil {
		return g
	}
	if len(g.noAppend) > 0 {
		g.root.addR(g.noAppend...)
		g.noAppend = []*route{}
	}
	return g
}

func (g *group) add(r *route) *group {
	if g.root == nil {
		g.noAppend = append(g.noAppend, r)
		return g
	}
	if len(g.noAppend) > 0 {
		g.root.addR(g.noAppend...)
		g.noAppend = []*route{}
	}
	g.root.addR(r)
	return g
}

func (g *group) tidy() *group {
	if g.root != nil {
		g.root.tidy()
	}
	return g
}

func (g *group) rtidy() bool {
	if g.root != nil {
		g.root.tidy()
		return true
	}
	return false
}
