package slim

import "strings"

type RouteShadow struct {
	Path       string
	Methods    []string
	Middleware []string
	Handler    string
	Name       string
	Label      []string
	Regexp     map[string]string
}

type RouteShadowSlice []RouteShadow

func (this RouteShadowSlice) Len() int {
	return len(this)
}

func (this RouteShadowSlice) Less(i, j int) bool {
	tmp := strings.Compare(this[i].Path, this[j].Path)
	if tmp < 0 {
		return true
	} else if tmp > 0 {
		return false
	}
	return strings.Compare(this[i].Handler, this[j].Handler) < 0
}

func (this RouteShadowSlice) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
