package slim

import (
	"net/url"
	"strconv"
)

func (this *Request) HasCookie(name string) bool {
	if cookie, err := this.r.Cookie(name); err == nil {
		if _, err = url.QueryUnescape(cookie.Value); err == nil {
			return true
		}
	}
	return false
}

func (this *Request) Cookie(name string, def ...string) string {
	if cookie, err := this.r.Cookie(name); err == nil {
		if r, err := url.QueryUnescape(cookie.Value); err == nil {
			return r
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

func (this *Request) CookieBool(name string, def ...bool) bool {
	if cookie, err := this.r.Cookie(name); err == nil {
		if data, err := url.QueryUnescape(cookie.Value); err == nil {
			if r, err := strconv.ParseBool(data); err == nil {
				return r
			}
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return false
}

func (this *Request) CookieInt(name string, def ...int) int {
	if cookie, err := this.r.Cookie(name); err == nil {
		if data, err := url.QueryUnescape(cookie.Value); err == nil {
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

func (this *Request) CookieFloat64(name string, def ...float64) float64 {
	if cookie, err := this.r.Cookie(name); err == nil {
		if data, err := url.QueryUnescape(cookie.Value); err == nil {
			if r, err := strconv.ParseFloat(data, 64); err == nil {
				return r
			}
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
