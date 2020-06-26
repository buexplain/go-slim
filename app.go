package slim

import (
	"github.com/buexplain/go-slim/tsmap"
	"github.com/buexplain/go-slim/view"
	"net/http"
	"strings"
	"sync"
)

type RecoverFunc func(ctx *Ctx, a interface{})

type ErrorFunc func(ctx *Ctx, err error)

type App struct {
	debug          bool
	//http请求的表单编码类型为multipart/form-data的内容解析到内存中的大小，超出会解析到磁盘
	formMaxMemory  int64
	//http请求的body的大小限制
	bodyMaxBytes   int64
	store          *tsmap.TSMap
	middleware     map[string][]Middleware
	pool           *sync.Pool
	mux            *Mux
	recoverFunc    RecoverFunc
	errorFunc      ErrorFunc
	view           *view.View
	sessionHandler SessionHandler
}

func New(debug bool) *App {
	tmp := new(App)
	tmp.debug = debug
	tmp.formMaxMemory = 10 << 20
	tmp.bodyMaxBytes = 10 << 20
	tmp.store = tsmap.New()
	tmp.middleware = map[string][]Middleware{
		http.MethodOptions: []Middleware{},
		http.MethodHead:    []Middleware{},
		http.MethodGet:     []Middleware{},
		http.MethodPost:    []Middleware{},
		http.MethodPatch:   []Middleware{},
		http.MethodPut:     []Middleware{},
		http.MethodDelete:  []Middleware{},
		http.MethodTrace:   []Middleware{},
		http.MethodConnect: []Middleware{},
	}
	tmp.pool = &sync.Pool{
		New: func() interface{} {
			return NewCtx(tmp, nil, nil)
		},
	}
	tmp.mux = NewMux()
	tmp.SetRecoverFunc(defaultRecoverFunc)
	tmp.SetErrorFunc(defaultErrorFunc)
	tmp.SetView(view.New("./view", !debug))
	return tmp
}

func (this *App) Debug() bool {
	return this.debug
}

func (this *App) SetFormMaxMemory(formMaxMemory int64) {
	this.formMaxMemory = formMaxMemory
}

func (this *App) SetBodyMaxBytes(bodyMaxBytes int64) {
	this.bodyMaxBytes = bodyMaxBytes
}

func (this *App) Store() *tsmap.TSMap {
	return this.store
}

func (this *App) Use(m Middleware, methods ...string) *App {
	if len(methods) == 0 || strings.ToUpper(methods[0]) == "ANY" {
		methods = []string{
			http.MethodOptions,
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPatch,
			http.MethodPut,
			http.MethodDelete,
			http.MethodTrace,
			http.MethodConnect,
		}
	}
	for _, v := range methods {
		v = strings.ToUpper(v)
		if _, ok := this.middleware[v]; !ok {
			panic("unknown http method")
		}
		this.middleware[v] = append(this.middleware[v], m)
	}
	return this
}

func (this *App) Mux() *Mux {
	return this.mux
}

func (this *App) SetRecoverFunc(recoverFunc RecoverFunc) {
	this.recoverFunc = recoverFunc
}

func (this *App) SetErrorFunc(errorFunc ErrorFunc) {
	this.errorFunc = errorFunc
}

func (this *App) SetView(view *view.View) {
	this.view = view
}

func (this *App) View() *view.View {
	return this.view
}

func (this *App) SetSessionHandler(sessionHandler SessionHandler) {
	this.sessionHandler = sessionHandler
}

func (this *App) SessionHandler() SessionHandler {
	return this.sessionHandler
}

func (this *App) Server(addr string) *http.Server {
	return &http.Server{Addr: addr, Handler: this}
}

func (this *App) Run(addr string) error {
	return this.Server(addr).ListenAndServe()
}

func (this *App) RunTLS(addr, certFile, keyFile string) error {
	return this.Server(addr).ListenAndServeTLS(certFile, keyFile)
}

func (this *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := this.pool.Get().(*Ctx)
	ctx.reset(w, r)
	defer func(app *App, context *Ctx) {
		if !context.w.send() {
			_ = context.w.send()
		}
		context.release()
		app.pool.Put(context)
	}(this, ctx)
	defer func(app *App, context *Ctx) {
		if a := recover(); a != nil {
			context.w.buffer.Reset()
			app.recoverFunc(context, a)
		}
	}(this, ctx)
	ctx.SetPath(r.URL.Path)
	ctx.Next()
}
