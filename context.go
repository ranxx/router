package router

import (
	"context"
	"reflect"
)

type cxtKey interface{}

const (
	ctxKeyCall   = "router.call"
	ctxKeyOrigin = "router.origin"
	ctxKeyExtra  = "router.extra"
)

// GetCallFromContext get Call from context
func GetCallFromContext(ctx context.Context) reflect.Value {
	v := ctx.Value(ctxKeyCall)
	if v == nil {
		return reflect.Value{}
	}
	return v.(reflect.Value)
}

// GetHandlerFromContext get Handler from context
func GetHandlerFromContext(ctx context.Context) Handler {
	v := ctx.Value(ctxKeyOrigin)
	if v == nil {
		return nil
	}
	return v.(Handler)
}

// GetFromMiddlewareContext get Middleware from context
func GetFromMiddlewareContext(ctx context.Context) Handler {
	v := ctx.Value(ctxKeyOrigin)
	if v == nil {
		return nil
	}
	return v.(Handler)
}

// GetExtraFromContext get Extra from context
func GetExtraFromContext(ctx context.Context) []Extra {
	v := ctx.Value(ctxKeyExtra)
	if v == nil {
		return nil
	}
	return v.([]Extra)
}

func newContext(m *metadata) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxKeyCall, m.call)
	ctx = context.WithValue(ctx, ctxKeyOrigin, m.origin)
	ctx = context.WithValue(ctx, ctxKeyExtra, m.extra)
	return ctx
}
