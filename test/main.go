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

func DiyRouter(md router.Customer[router.Handler], args ...interface{}) bool {
	fmt.Println("router begin", md.Call().String(), args)
	ins := make([]reflect.Value, 0, len(args))
	for _, v := range args {
		ins = append(ins, reflect.ValueOf(v))
	}
	ctx := args[0].(*Context)
	fmt.Printf("ctx: %v\n", ctx)
	md.Call().Call(ins)
	if ctx.Abort {
		return false
	}
	return true
}

func main() {
	fmt.Println("Hello World")
	r := router.NewRouter(router.WithRouter(DiyRouter), router.WithMiddleware(func(em router.Customer[router.Middleware], args ...interface{}) bool {
		fmt.Println("skip all mid", em.Extra(), em.Call().String())
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

	g := r.NewG(func(ctx *Context, req interface{}) string {
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
