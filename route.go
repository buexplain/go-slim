package slim

import (
	"net/http"
	"regexp"
	"strings"
)

type RouteSetInterface interface {
	SetName(name string) RouteSetInterface
	AddLabel(label ...string) RouteSetInterface
	Use(m ...Middleware) RouteSetInterface
	Regexp(key string, pattern string) RouteSetInterface
}

type RouteGetInterface interface {
	GetPath() string
	GetName() string
	HasLabel(label string) bool
	GetLabel() []string
}

type Route struct {
	mux        *Mux
	path       string
	methods    []string
	middleware []Middleware
	handler    Handler
	name       string
	label      []string
	regexp     map[string]*regexp.Regexp
}

func NewRoute(mux *Mux, path string, methods []string, handler Handler) *Route {
	tmp := new(Route)
	tmp.mux = mux
	tmp.setPath(path)
	tmp.methods = []string{}
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
	for _, method := range methods {
		method = strings.ToUpper(method)
		if method == http.MethodOptions ||
			method == http.MethodHead ||
			method == http.MethodGet ||
			method == http.MethodPost ||
			method == http.MethodPatch ||
			method == http.MethodPut ||
			method == http.MethodDelete ||
			method == http.MethodTrace ||
			method == http.MethodConnect {
			tmp.methods = append(tmp.methods, method)
		} else {
			panic("unknown http method: " + method)
		}
	}
	if len(tmp.methods) == 0 {
		panic("unknown http method")
	}
	tmp.middleware = []Middleware{}
	tmp.handler = handler
	tmp.label = []string{}
	tmp.regexp = make(map[string]*regexp.Regexp, len(mux.regexp))
	for k, v := range mux.regexp {
		tmp.regexp[k] = v
	}
	return tmp
}

func (this *Route) GetPath() string {
	return this.path
}

func (this *Route) setPath(path string) {
	this.path = "/" + strings.Trim(path, "/")
}

func (this *Route) SetName(name string) RouteSetInterface {
	if _, ok := this.mux.routeMap[name]; ok {
		panic("route name already exists: " + name)
	}
	delete(this.mux.routeMap, this.name)
	this.name = name
	this.mux.routeMap[name] = this
	return this
}

func (this *Route) GetName() string {
	return this.name
}

func (this *Route) AddLabel(label ...string) RouteSetInterface {
	for _, v := range label {
		if v == "" || this.HasLabel(v) {
			continue
		}
		this.label = append(this.label, v)
	}
	return this
}

func (this *Route) HasLabel(label string) bool {
	for _, v := range this.label {
		if v == label {
			return true
		}
	}
	return false
}

func (this *Route) GetLabel() []string {
	return this.label
}

func (this *Route) Use(m ...Middleware) RouteSetInterface {
	for _, v := range m {
		if v == nil {
			continue
		}
		this.middleware = append(this.middleware, v)
	}
	return this
}

func (this *Route) Regexp(key string, pattern string) RouteSetInterface {
	this.regexp[key] = regexp.MustCompile(pattern)
	return this
}
