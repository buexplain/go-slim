package slim

import (
	"bytes"
	"fmt"
	"github.com/buexplain/go-slim/tree"
	"github.com/olekukonko/tablewriter"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type Mux struct {
	data         map[string]*tree.Tree
	routeMap     map[string]*Route
	groups       []*RouteGroup
	regexp       map[string]*regexp.Regexp
	defaultRoute *Route
}

func NewMux() *Mux {
	tmp := new(Mux)
	tmp.data = map[string]*tree.Tree{
		http.MethodOptions: tree.New(false),
		http.MethodHead:    tree.New(false),
		http.MethodGet:     tree.New(false),
		http.MethodPost:    tree.New(false),
		http.MethodPatch:   tree.New(false),
		http.MethodPut:     tree.New(false),
		http.MethodDelete:  tree.New(false),
		http.MethodTrace:   tree.New(false),
		http.MethodConnect: tree.New(false),
	}
	tmp.routeMap = make(map[string]*Route)
	tmp.groups = make([]*RouteGroup, 0)
	tmp.regexp = make(map[string]*regexp.Regexp)
	tmp.SetDefaultRoute(defaultRoute)
	return tmp
}

func (this *Mux) AddRoute(route *Route) error {
	if len(this.groups) > 0 {
		route.setPath(this.groups[len(this.groups)-1:][0].prefix + route.path)
	}
	for _, method := range route.methods {
		if err := this.data[method].Add(route.path, route); err != nil {
			return fmt.Errorf("method: %s path: %s err: %+v", method, route.path, err)
		}
	}
	if this.groups != nil {
		for _, v := range this.groups {
			v.addRoute(route)
		}
	}
	route.SetName(strconv.Itoa(len(this.routeMap)))
	return nil
}

func (this *Mux) Wrap(handler interface{}) Handler {
	switch h := handler.(type) {
	case Handler:
		return h
	case func(ctx *Ctx, w *Response, r *Request) error:
		return h
	case http.HandlerFunc, http.Handler:
		return func(ctx *Ctx, w *Response, r *Request) error {
			h.(http.Handler).ServeHTTP(w, r.Raw())
			return nil
		}
	case func(http.ResponseWriter, *http.Request):
		return func(ctx *Ctx, w *Response, r *Request) error {
			h(w, r.Raw())
			return nil
		}
	default:
		panic("unknown mux handler")
	}
}

func (this *Mux) Add(path string, handler Handler, method string) RouteSetInterface {
	route := NewRoute(this, path, []string{method}, handler)
	if err := this.AddRoute(route); err != nil {
		panic(err)
	}
	return route
}

func (this *Mux) Options(path string, handler Handler) RouteSetInterface {
	return this.Add(path, handler, http.MethodOptions)
}

func (this *Mux) Head(path string, handler Handler) RouteSetInterface {
	return this.Add(path, handler, http.MethodHead)
}

func (this *Mux) Get(path string, handler Handler) RouteSetInterface {
	return this.Add(path, handler, http.MethodGet)
}

func (this *Mux) Post(path string, handler Handler) RouteSetInterface {
	return this.Add(path, handler, http.MethodPost)
}

func (this *Mux) Patch(path string, handler Handler) RouteSetInterface {
	return this.Add(path, handler, http.MethodPatch)
}

func (this *Mux) Put(path string, handler Handler) RouteSetInterface {
	return this.Add(path, handler, http.MethodPut)
}

func (this *Mux) Delete(path string, handler Handler) RouteSetInterface {
	return this.Add(path, handler, http.MethodDelete)
}

func (this *Mux) Trace(path string, handler Handler) RouteSetInterface {
	return this.Add(path, handler, http.MethodTrace)
}

func (this *Mux) Connect(path string, handler Handler) RouteSetInterface {
	return this.Add(path, handler, http.MethodConnect)
}

func (this *Mux) Any(path string, handler Handler, methods ...string) RouteSetInterface {
	var route *Route = NewRoute(this, path, methods, handler)
	if err := this.AddRoute(route); err != nil {
		panic(err)
	}
	return route
}

func (this *Mux) Restful(path string, handler ...Handler) RouteSetInterface {
	routeRestful := &RouteRestful{data: []*Route{}}
	for _, h := range handler {
		fn := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
		tmp := strings.LastIndex(fn, ".")
		if tmp == -1 {
			continue
		}
		fn = fn[tmp+1:]
		switch true {
		case strings.EqualFold(fn, "index"):
			route := this.Get(path, h)
			routeRestful.data = append(routeRestful.data, route.(*Route))
			break
		case strings.EqualFold(fn, "create"):
			route := this.Get(strings.TrimRight(path, "/")+"/create", h)
			routeRestful.data = append(routeRestful.data, route.(*Route))
			break
		case strings.EqualFold(fn, "store"):
			route := this.Post(path, h)
			routeRestful.data = append(routeRestful.data, route.(*Route))
			break
		case strings.EqualFold(fn, "edit"):
			route := this.Get(strings.TrimRight(path, "/")+"/edit/:id", h)
			routeRestful.data = append(routeRestful.data, route.(*Route))
			break
		case strings.EqualFold(fn, "update"):
			route := this.Put(strings.TrimRight(path, "/")+"/update/:id", h)
			routeRestful.data = append(routeRestful.data, route.(*Route))
			break
		case strings.EqualFold(fn, "destroy"):
			route := this.Delete(strings.TrimRight(path, "/")+"/delete/:id", h)
			routeRestful.data = append(routeRestful.data, route.(*Route))
			break
		case strings.EqualFold(fn, "show"):
			route := this.Get(strings.TrimRight(path, "/")+"/show/:id", h)
			routeRestful.data = append(routeRestful.data, route.(*Route))
			break
		}
	}
	return routeRestful
}

func (this *Mux) Group(prefix string, f func()) *RouteGroup {
	if len(this.groups) > 0 {
		if prefix != "" && prefix != "/" {
			prefix = this.groups[len(this.groups)-1:][0].prefix + "/" + strings.Trim(prefix, "/")
		} else {
			prefix = this.groups[len(this.groups)-1:][0].prefix
		}
	}
	g := NewRouteGroup(prefix)
	this.groups = append(this.groups, g)
	f()
	this.groups = this.groups[:len(this.groups)-1]
	return g
}

func (this *Mux) Regexp(key string, pattern string) *Mux {
	this.regexp[key] = regexp.MustCompile(pattern)
	return this
}

func (this *Mux) SetDefaultRoute(handler Handler) RouteSetInterface {
	this.defaultRoute = NewRoute(this, "", nil, handler)
	return this.defaultRoute
}

func (this *Mux) GetRouteByName(name string) RouteGetInterface {
	if r, ok := this.routeMap[name]; ok {
		return r
	}
	return nil
}

func (this *Mux) GetRouteMap(packageName ...string) RouteShadowSlice {
	if len(packageName) == 0 {
		packageName = append(packageName, "")
	}
	shadows := make(RouteShadowSlice, 0, len(this.routeMap))
	for _, route := range this.routeMap {
		shadow := RouteShadow{}
		shadow.Path = route.path
		shadow.Methods = route.methods
		shadow.Middleware = make([]string, 0, len(route.middleware))
		shadow.Handler = strings.TrimPrefix(runtime.FuncForPC(reflect.ValueOf(route.handler).Pointer()).Name(), packageName[0])
		shadow.Name = route.name
		shadow.Label = route.label
		for _, m := range route.middleware {
			shadow.Middleware = append(shadow.Middleware, strings.TrimPrefix(runtime.FuncForPC(reflect.ValueOf(m).Pointer()).Name(), packageName[0]))
		}
		shadow.Regexp = map[string]string{}
		for k, v := range route.regexp {
			shadow.Regexp[k] = v.String()
		}
		shadows = append(shadows, shadow)
	}
	sort.Sort(shadows)
	return shadows
}

func (this *Mux) DumpRouteMap(packageName ...string) string {
	buf := bytes.NewBuffer(nil)
	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"#", "Path", "Methods", "Middleware", "Handler", "Name", "Label", "Regexp"})
	table.SetRowLine(true)
	table.SetBorder(true)
	table.SetAutoWrapText(false)
	shadows := this.GetRouteMap(packageName...)
	for k, v := range shadows {
		reg := bytes.Buffer{}
		for k2, v2 := range v.Regexp {
			if reg.Len() > 0 {
				reg.WriteString("\n")
			}
			reg.WriteString(k2)
			reg.WriteString(":")
			reg.WriteString(v2)
		}
		var methods string = "ANY"
		if len(v.Methods) != 9 {
			methods = strings.Join(v.Methods, "\n")
		}
		tmp := []string{strconv.Itoa(k + 1), v.Path, methods, strings.Join(v.Middleware, "\n"), v.Handler, v.Name, strings.Join(v.Label, "\n"), reg.String()}
		table.Append(tmp)
	}
	table.Render()
	return buf.String()
}

func (this *Mux) match(ctx *Ctx) *Route {
	if currTree, ok := this.data[ctx.r.r.Method]; ok {
		if result, ok := currTree.Search(ctx.Path(), ctx.r.param); ok {
			if route, ok := result.(*Route); ok {
				for k, v := range route.regexp {
					if ctx.r.HasParam(k) {
						if v.MatchString(ctx.r.Param(k)) == false {
							return this.defaultRoute
						}
					}
				}
				return route
			} else {
				return this.defaultRoute
			}
		} else {
			return this.defaultRoute
		}
	} else {
		return this.defaultRoute
	}
}
