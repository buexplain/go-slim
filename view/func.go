package view

import (
	"fmt"
	"html/template"
)

func HTML(a interface{}) template.HTML {
	if s, ok := a.(string); ok {
		return template.HTML(s)
	}
	if e, ok := a.(error); ok {
		return template.HTML(e.Error())
	}
	return template.HTML(fmt.Sprintf("%+v", a))
}
