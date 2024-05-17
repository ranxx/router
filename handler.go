package router

// network
// Context interface{}

// type Client struct {
// }

// type Context struct {
// 	context.Context
// 	c *Client
// }

// type XXXRequest struct {
// }

// func Htest(ctx Context, req XXXRequest) {
// 	// mids
// 	// ctx Context,

// 	ctx = context.WithValue(ctx, "", "")
// }

// // Handler ...
// type Handler interface {
// 	Serve(context.Context, *request.Request)
// }

// // WrapHandler ...
// type WrapHandler func(context.Context, *request.Request)

// // Serve ...
// func (w WrapHandler) Serve(ctx context.Context, req *request.Request) {
// 	w(ctx, req)
// }

// // Request request
// type Request struct {
// 	M     message.Messager
// 	C     conner.Conner
// 	abort bool
// }

// // NewRequest new requests
// func NewRequest(m message.Messager, c conner.Conner) *Request {
// 	return &Request{
// 		M: m,
// 		C: c,
// 	}
// }

// // Abort ...
// func (r *Request) Abort() {
// 	r.abort = true
// }

// // GetAbort ...
// func (r *Request) GetAbort() bool {
// 	return r.abort
// }
