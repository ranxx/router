package router

import "reflect"

// RouterType router type
type RouterType int

// const
const (
	Type  RouterType = 0
	Value RouterType = 1
)

// Handler is a function type. If the Opts Router Exec is not set,
// then the function parameter of the Handler must be the same as the number of ARGS parameters of Router
type Handler interface{}

// Middleware is a function type. If the Opts Middleware Exec is not set,
// then the function parameter of the Middleware must be the same as the number of ARGS parameters of Router
type Middleware interface{}

// Path path
type Path interface{}

// Extra extra
type Extra interface{}

type Root interface {
	// NewGroup new group, Middleware is a function type
	NewGroup(mids ...Middleware) Group

	// NewRouter new router, handler is a function type
	NewRouter(path Path, handler Handler, handlers ...Handler) Router

	// Use same as Middleware func, Middleware is a function type
	Use(mids ...Middleware) Root

	// Add same as Register func, Handler is a function type
	Add(path Path, handler Handler, handlers ...Handler) Root

	// Middleware same as Use func, Middleware is a function type
	Middleware(mids ...Middleware) Root

	// Register same as Add func, Handler is a function type
	Register(path Path, handler Handler, extra ...Extra) Root

	// Router router
	Router(path Path, args ...interface{})

	// Handlers get handlers
	Handlers(path Path) []reflect.Value

	// Range fc(Path, []Handler, [extra...])
	Range(fc interface{})
}

type Group interface {
	// Root return Root interface
	Root() Root

	// NewRouter new Router, Handler is a function type
	NewRouter(path Path, handler Handler, handlers ...Handler) Router

	// Use same as Middleware func, Middleware is a function type
	Use(mids ...Middleware) Group

	// Add same as Register func, Handler is a function type
	Add(path Path, handler Handler, handlers ...Handler) Group

	// Middleware same as Use func, Middleware is a function type
	Middleware(mids ...Middleware) Group

	// Register same as Add func, Handler is a function type
	Register(path Path, handler Handler, extra ...Extra) Group

	// WithExtra with extra
	WithExtra(extra ...Extra) Group
}

type Router interface {
	// WithExtra with extra
	WithExtra(extra ...Extra) Router
}
