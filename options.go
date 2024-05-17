package router

import (
	"fmt"
	"reflect"
)

type Labels map[string]interface{}

type IHandler interface {
	Handler | Middleware
}

type Customer[T IHandler] interface {
	Call() reflect.Value
	Origin() T
	Extra() []Extra
}

// type ExecMetadata[T IHandler] struct {
// 	Call   reflect.Value
// 	Origin T
// 	Extra  []Extra
// }

type Exec[T IHandler] func(p Customer[T], args ...interface{}) bool

// // CustomeRouter len(args) == len(Handler args) == len(Middleware args) == len(Router args)
// type CustomeRouter func(orgin Handler, call reflect.Value, args ...interface{}) bool

// // CustomeMiddleware len(args) == len(Handler args) == len(Middleware args) == len(Router args)
// type CustomeMiddleware func(origin Middleware, call reflect.Value, args ...interface{}) bool

// Options options
type Options struct {
	Type       RouterType       `default:"0" md:"router type"`
	Router     Exec[Handler]    `md:"return false break"`
	Middleware Exec[Middleware] `md:"return false break"`
}

// Option option
type Option func(*Options)

// WithType with route type
func WithType(rt RouterType) Option {
	return func(o *Options) {
		o.Type = rt
	}
}

// WithRouter with custome router exec
func WithRouter(router Exec[Handler]) Option {
	return func(o *Options) {
		o.Router = router
	}
}

// WithMiddleware with custome middleware exec
func WithMiddleware(router Exec[Middleware]) Option {
	return func(o *Options) {
		o.Middleware = router
	}
}

func newDefaultOptions() *Options {
	return &Options{
		Type: Type,
	}
}

func checkHandler(handlers ...Handler) []Handler {
	for _, h := range handlers {
		val := reflect.ValueOf(h)
		if val.Type().Kind() != reflect.Func {
			panic(fmt.Sprintf("Router Handler %s must be func", val.String()))
		}
	}
	return handlers
}

func checkMiddleware(mids ...Middleware) []Middleware {
	for _, h := range mids {
		val := reflect.ValueOf(h)
		if val.Type().Kind() != reflect.Func {
			panic(fmt.Sprintf("Router Middleware %s must be func", val.String()))
		}
	}
	return mids
}
