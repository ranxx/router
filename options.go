package router

import (
	"context"
	"fmt"
	"reflect"
)

// Labels labels
type Labels map[string]interface{}

// IHandler func
type IHandler interface {
	Handler | Middleware
}

// Exec custome Handler Middleware exec
type Exec func(ctx context.Context, args ...interface{}) bool

// Options options
type Options struct {
	ty         RouterType `default:"0" md:"router type"`
	router     Exec       `md:"return false break"`
	middleware Exec       `md:"return false break"`
}

// Option option
type Option func(*Options)

// WithType with route type
func WithType(rt RouterType) Option {
	return func(o *Options) {
		o.ty = rt
	}
}

// WithRouter with custome router exec
func WithRouter(router Exec) Option {
	return func(o *Options) {
		o.router = router
	}
}

// WithMiddleware with custome middleware exec
func WithMiddleware(router Exec) Option {
	return func(o *Options) {
		o.middleware = router
	}
}

func newDefaultOptions() *Options {
	return &Options{
		ty: Type,
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
