package slim

import (
	"github.com/buexplain/go-slim/tsmap"
	"net/http"
)

//请求上下文
type Ctx struct {
	//当前的app
	app *App
	//上下文存储容器
	store *tsmap.TSMap
	//当前响应对象
	w *Response
	//当前请求对象
	r *Request
	//中间件循环所需变量
	nextI int
	//当前请求的路由中间件
	middleware []Middleware
	//中间件循环所需变量
	nextJ int
	//当前请求命中的路由
	route *Route
	//用于路由匹配的path
	routeMatchPath string
}

//新建一个上下文
func NewCtx(app *App, w http.ResponseWriter, r *http.Request) *Ctx {
	tmp := &Ctx{
		app:        app,
		store:      tsmap.New(),
		nextI:      0,
		middleware: nil,
		nextJ:      0,
		route:      nil,
	}
	tmp.r = NewRequest(tmp, r)
	tmp.w = NewResponse(tmp, w)
	return tmp
}

//重置上下文相关字段
func (this *Ctx) reset(w http.ResponseWriter, r *http.Request) {
	this.w.Reset(w)
	this.r.Reset(r)
	if r != nil {
		this.middleware = this.app.middleware[r.Method]
	}
	this.route = nil
}

//释放上下文相关字段
func (this *Ctx) release() {
	this.w.release()
	this.r.release()
	this.store.Release()
	this.nextI = 0
	this.middleware = nil
	this.nextJ = 0
	this.route = nil
	this.routeMatchPath = ""
}

//返回上下文存储容器
func (this *Ctx) Store() *tsmap.TSMap {
	return this.store
}

//返回当前的响应对象
func (this *Ctx) Response() *Response {
	return this.w
}

//返回当前的请求对象
func (this *Ctx) Request() *Request {
	return this.r
}

//返回当前请求命中的路由，此方法必须在全局中间件结束后调用，否则返回nil
func (this *Ctx) Route() RouteGetInterface {
	if this.route == nil {
		return nil
	}
	return this.route
}

//进入下一个中间件
func (this *Ctx) Next() {
	if this.nextI < len(this.middleware) {
		this.nextI++
		this.middleware[this.nextI-1](this, this.w, this.r)
	} else {
		if this.route == nil {
			this.route = this.app.mux.match(this)
		}
		if this.nextJ < len(this.route.middleware) {
			this.nextJ++
			this.route.middleware[this.nextJ-1](this, this.w, this.r)
		} else {
			if this.nextJ == len(this.route.middleware) {
				this.nextJ++
				this.Throw(this.route.handler(this, this.w, this.r))
			}
		}
	}
}

//跳出全局或路由中间件
func (this *Ctx) Break() {
	if this.route == nil {
		//未命中路由，跳出全局中间件
		this.nextI = len(this.middleware)
		this.Next()
	} else {
		//已经命中路由，跳出路由中间件
		if this.nextJ < len(this.route.middleware) {
			this.nextJ = len(this.route.middleware)
			this.Next()
		}
	}
}

//抛出一个错误
func (this *Ctx) Throw(err error) {
	this.app.errorFunc(this, err)
}

//返回app
func (this *Ctx) App() *App {
	return this.app
}

//设置用于路由匹配的path
func (this *Ctx) SetPath(path string) {
	l := len(path)
	c := "/"
	if l == 0 {
		path = c
	} else {
		if path[0:1] != c {
			path = c + path
		}
		l = len(path)
		if l > 1 && path[l-1:] == c {
			path = path[:l-1]
		}
	}
	this.routeMatchPath = path
}

//获取用于路由匹配的path
func (this *Ctx) Path() string {
	return this.routeMatchPath
}
