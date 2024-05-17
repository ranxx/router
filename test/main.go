package main

import (
	"context"
	"fmt"
	"reflect"

	router "github.com/ranxx/router"
)

type Context struct {
	context.Context
	Abort bool
	Name  string
}

func DiyRouter(ctx context.Context, args ...interface{}) bool {
	fmt.Println("router begin", router.GetCallFromContext(ctx).String(), args)
	ins := make([]reflect.Value, 0, len(args))
	for _, v := range args {
		ins = append(ins, reflect.ValueOf(v))
	}
	rctx := args[0].(*Context)
	fmt.Printf("router rctx: %v\n", rctx)
	router.GetCallFromContext(ctx).Call(ins)
	if rctx.Abort {
		return false
	}
	return true
}

func main() {
	fmt.Println("Hello World")
	r := router.NewRouter(router.WithRouter(DiyRouter), router.WithMiddleware(func(ctx context.Context, args ...interface{}) bool {
		fmt.Println("skip all mid", router.GetExtraFromContext(ctx), router.GetCallFromContext(ctx).String())
		return true
	}))

	r.Use(func(ctx *Context, req interface{}) {
		fmt.Println("mid call")
		ctx.Name = "mid call success"
	})

	r.Add("x", func(ctx *Context, req string) {
		fmt.Printf("string ctx: %v req: %v\n", ctx, req)
		ctx.Abort = true
	}, func(ctx *Context, req string) {
		fmt.Printf("string-2 ctx: %v req: %v\n", ctx, req)
	})

	g := r.NewGroup(func(ctx *Context, req interface{}) string {
		fmt.Println("group 1")
		return ""
	}).WithExtra(router.Labels{"name": "group 1"})

	g.Add(0, func(ctx *Context, req string) int {
		fmt.Printf("int ctx: %v req: %v\n", ctx, req)
		return 99
	})
	// fmt.Println(r.Handlers(0))
	r.Router("string", &Context{Context: context.Background()}, "lili")
	fmt.Println("---------------------------------------------")
	r.Router(0, &Context{Context: context.Background()}, "")
}
