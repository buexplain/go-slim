package slim

import (
	"net/url"
	"strconv"
)

func (this *Request) ParseQuery() error {
	if this.query == nil {
		if this.r.URL == nil {
			this.query = make(url.Values)
		} else {
			var err error
			this.query, err = url.ParseQuery(this.r.URL.RawQuery)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *Request) hasQuery(key string) (string, bool) {
	err := this.ParseQuery()
	if err != nil {
		//解析出错
		return key, false
	}
	//判断是否存在
	_, ok := this.query[key]
	if !ok {
		//判断是否存在带中括号类型的
		l := len(key)
		if l <= 2 || key[l-2:] != "[]" {
			tmp := key + "[]"
			_, ok = this.query[tmp]
			if ok {
				//存在则将key重置
				key = tmp
			}
		}
	}
	return key, ok
}

func (this *Request) HasQuery(key string) bool {
	_, ok := this.hasQuery(key)
	return ok
}

func (this *Request) SetQuery(key string, value string) {
	key, _ = this.hasQuery(key)
	this.query.Set(key, value)
}

func (this *Request) AddQuery(key string, value string) {
	key, _ = this.hasQuery(key)
	this.query.Add(key, value)
}

func (this *Request) DelQuery(key string) {
	key, _ = this.hasQuery(key)
	this.query.Del(key)
}

func (this *Request) EncodeQuery() string {
	_ = this.ParseQuery()
	return this.query.Encode()
}

func (this *Request) Query(key string, def ...string) string {
	var ok bool
	key, ok = this.hasQuery(key)
	if ok {
		data := this.query[key]
		if len(data) > 0 {
			return data[0]
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

func (this *Request) QueryBool(key string, def ...bool) bool {
	var ok bool
	key, ok = this.hasQuery(key)
	if ok {
		data := this.query[key]
		if len(data) > 0 {
			if r, err := strconv.ParseBool(data[0]); err == nil {
				return r
			}
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return false
}

func (this *Request) QueryInt(key string, def ...int) int {
	var ok bool
	key, ok = this.hasQuery(key)
	if ok {
		data := this.query[key]
		if len(data) > 0 {
			if r, err := strconv.Atoi(data[0]); err == nil {
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
func (this *Request) QueryPositiveInt(key string, def ...int) int {
	var ok bool
	key, ok = this.hasQuery(key)
	if ok {
		data := this.query[key]
		if len(data) > 0 {
			if r, err := strconv.Atoi(data[0]); err == nil && r > 0 {
				return r
			}
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (this *Request) QueryFloat64(key string, def ...float64) float64 {
	var ok bool
	key, ok = this.hasQuery(key)
	if ok {
		data := this.query[key]
		if len(data) > 0 {
			if r, err := strconv.ParseFloat(data[0], 64); err == nil {
				return r
			}
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (this *Request) QuerySlice(key string, def ...[]string) []string {
	var ok bool
	key, ok = this.hasQuery(key)
	if ok {
		return this.query[key]
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

func (this *Request) QuerySliceBool(key string, def ...[]bool) []bool {
	data := this.QuerySlice(key)
	if data != nil {
		result := make([]bool, 0, len(data))
		for _, v := range data {
			if r, err := strconv.ParseBool(v); err == nil {
				result = append(result, r)
			}
		}
		return result
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

func (this *Request) QuerySliceInt(key string, def ...[]int) []int {
	data := this.QuerySlice(key)
	if data != nil {
		result := make([]int, 0, len(data))
		for _, v := range data {
			if r, err := strconv.Atoi(v); err == nil {
				result = append(result, r)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

func (this *Request) QuerySlicePositiveInt(key string, def ...[]int) []int {
	data := this.QuerySlice(key)
	if data != nil {
		result := make([]int, 0, len(data))
		for _, v := range data {
			if r, err := strconv.Atoi(v); err == nil && r > 0 {
				result = append(result, r)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

func (this *Request) QuerySliceFloat64(key string, def ...[]float64) []float64 {
	data := this.QuerySlice(key)
	if data != nil {
		result := make([]float64, 0, len(data))
		for _, v := range data {
			if r, err := strconv.ParseFloat(v, 64); err == nil {
				result = append(result, r)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

//获取符合此类 (.+?)\\[(.+?)\\] 正则的参数
func (this *Request) QueryMap(key string, def ...map[string]string) map[string]string {
	err := this.ParseQuery()
	if err != nil {
		if len(def) > 0 {
			return def[0]
		}
		return nil
	}
	key_l := len(key)
	var data map[string]string
	var k_l int
	for k, v := range this.query {
		k_l = len(k)
		if len(v) > 0 && k_l-2 > key_l && k[:key_l] == key && k[key_l:key_l+1] == "[" && k[k_l-1:] == "]" {
			if data == nil {
				data = map[string]string{}
			}
			data[k[key_l+1:k_l-1]] = v[0]
		}
	}
	if data != nil {
		return data
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}
