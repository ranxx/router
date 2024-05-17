package main

import (
	"context"
	"fmt"
	"reflect"

	"github.com/ranxx/router"
)

type rcontext struct {
	context.Context
}

/*
type Labels map[string]interface{}

type IHandler interface {
	Handler | Middleware
}

type Customer[T IHandler] interface {
	Call() reflect.Value
	Origin() T
	Extra() []Extra
}
*/

type CxtKey interface{}

func GetFromContext[T interface{}](ctx context.Context, key CxtKey) T {
	return ctx.Value(key).(T)
}

func main() {
	var h interface{}
	h = func() {
		fmt.Println("called")
	}
	ctx := context.WithValue(context.TODO(), "ccc", h)

	h2 := GetFromContext[router.IHandler](ctx, "ccc")

	fmt.Println(h2)

	reflect.ValueOf(h2).Call([]reflect.Value{})
}
