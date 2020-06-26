package slim

import (
	"reflect"
	"runtime"
	"strings"
)

type RouteRestful struct {
	data []*Route
}

func (this *RouteRestful) SetName(name string) RouteSetInterface {
	for _, v := range this.data {
		fn := runtime.FuncForPC(reflect.ValueOf(v.handler).Pointer()).Name()
		tmp := strings.LastIndex(fn, ".")
		if tmp == -1 {
			continue
		}
		fn = fn[tmp+1:]
		switch true {
		case strings.EqualFold(fn, "index"):
			v.SetName(name + ".index")
			break
		case strings.EqualFold(fn, "create"):
			v.SetName(name + ".create")
			break
		case strings.EqualFold(fn, "store"):
			v.SetName(name + ".store")
			break
		case strings.EqualFold(fn, "edit"):
			v.SetName(name + ".edit")
			break
		case strings.EqualFold(fn, "update"):
			v.SetName(name + ".update")
			break
		case strings.EqualFold(fn, "destroy"):
			v.SetName(name + ".destroy")
			break
		case strings.EqualFold(fn, "show"):
			v.SetName(name + ".show")
			break
		}
	}
	return this
}

func (this *RouteRestful) AddLabel(label ...string) RouteSetInterface {
	for _, v := range this.data {
		v.AddLabel(label...)
	}
	return this
}

func (this *RouteRestful) Use(m ...Middleware) RouteSetInterface {
	for _, v := range this.data {
		v.Use(m...)
	}
	return this
}

func (this *RouteRestful) Regexp(key string, pattern string) RouteSetInterface {
	for _, v := range this.data {
		v.Regexp(key, pattern)
	}
	return this
}
