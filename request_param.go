package slim

import "strconv"

func (this *Request) HasParam(key string) bool {
	return this.param.Has(key)
}

func (this *Request) Param(key string, def ...string) string {
	if this.param.Has(key) {
		if data, ok := this.param.Get(key).(string); ok {
			return data
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

func (this *Request) ParamBool(key string, def ...bool) bool {
	if this.param.Has(key) {
		if data, ok := this.param.Get(key).(string); ok {
			if data, err := strconv.ParseBool(data); err == nil {
				return data
			}
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return false
}

func (this *Request) ParamInt(key string, def ...int) int {
	if this.param.Has(key) {
		if data, ok := this.param.Get(key).(string); ok {
			if r, err := strconv.Atoi(data); err == nil {
				return r
			}
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

//返回正整数
func (this *Request) ParamPositiveInt(key string, def ...int) int {
	if this.param.Has(key) {
		if data, ok := this.param.Get(key).(string); ok {
			if r, err := strconv.Atoi(data); err == nil && r > 0 {
				return r
			}
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (this *Request) ParamFloat64(key string, def ...float64) float64 {
	if this.param.Has(key) {
		if data, ok := this.param.Get(key).(string); ok {
			if n, err := strconv.ParseFloat(data, 64); err == nil {
				return n
			}
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
