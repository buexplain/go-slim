package slim

import (
	"fmt"
	"regexp"
	"strings"
)

type RouteGroup struct {
	prefix     string
	data       []*Route
	middleware map[string][]Middleware
}

func NewRouteGroup(prefix string) *RouteGroup {
	if prefix == "/" {
		prefix = ""
	}
	if prefix != "" {
		prefix = "/" + strings.Trim(prefix, "/")
	}
	return &RouteGroup{prefix: prefix, data: make([]*Route, 0), middleware: map[string][]Middleware{}}
}

func (this *RouteGroup) addRoute(route *Route) {
	this.data = append(this.data, route)
}

func (this *RouteGroup) AddLabel(label ...string) *RouteGroup {
	for _, v := range this.data {
		v.AddLabel(label...)
	}
	return this
}

func (this *RouteGroup) Use(m ...Middleware) *RouteGroup {
	if m == nil {
		return this
	}
	for _, v := range this.data {
		key := fmt.Sprintf("%p", v)
		if _, ok := this.middleware[key]; !ok {
			this.middleware[key] = []Middleware{}
		}
		self := v.middleware[len(this.middleware[key]):]
		v.middleware = []Middleware{}
		this.middleware[key] = append(this.middleware[key], m...)
		v.Use(this.middleware[key]...)
		v.Use(self...)
	}
	return this
}

func (this *RouteGroup) Regexp(key string, pattern string) *RouteGroup {
	tmp := regexp.MustCompile(pattern)
	for _, v := range this.data {
		if _, ok := v.regexp[key]; !ok {
			v.regexp[key] = tmp
		}
	}
	return this
}
