package slim

type Middleware func(ctx *Ctx, w *Response, r *Request)
