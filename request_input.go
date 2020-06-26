package slim

func (this *Request) HasInput(key string) bool {
	if this.HasParam(key) {
		return true
	}
	if this.HasQuery(key) {
		return true
	}
	if this.HasForm(key) {
		return true
	}
	return true
}

func (this *Request) Input(key string, def ...string) string {
	switch true {
	case this.HasParam(key):
		return this.Param(key, def...)
	case this.HasQuery(key):
		return this.Query(key, def...)
	case this.HasForm(key):
		return this.Form(key, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return ""
	}
}

func (this *Request) InputBool(key string, def ...bool) bool {
	switch true {
	case this.HasParam(key):
		return this.ParamBool(key, def...)
	case this.HasQuery(key):
		return this.QueryBool(key, def...)
	case this.HasForm(key):
		return this.FormBool(key, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return false
	}
}

func (this *Request) InputInt(key string, def ...int) int {
	switch true {
	case this.HasParam(key):
		return this.ParamInt(key, def...)
	case this.HasQuery(key):
		return this.QueryInt(key, def...)
	case this.HasForm(key):
		return this.FormInt(key, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}

func (this *Request) InputPositiveInt(key string, def ...int) int {
	switch true {
	case this.HasParam(key):
		return this.ParamPositiveInt(key, def...)
	case this.HasQuery(key):
		return this.QueryPositiveInt(key, def...)
	case this.HasForm(key):
		return this.FormPositiveInt(key, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}

func (this *Request) InputFloat64(key string, def ...float64) float64 {
	switch true {
	case this.HasParam(key):
		return this.ParamFloat64(key, def...)
	case this.HasQuery(key):
		return this.QueryFloat64(key, def...)
	case this.HasForm(key):
		return this.FormFloat64(key, def...)
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}
