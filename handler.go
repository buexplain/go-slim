package slim

type Handler func(ctx *Ctx, w *Response, r *Request) error
