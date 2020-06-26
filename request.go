package slim

import (
	"github.com/buexplain/go-slim/constant"
	"github.com/buexplain/go-slim/tsmap"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	ctx     *Ctx
	r       *http.Request
	query   url.Values
	param   *tsmap.TSMap
	session Session
}

func NewRequest(ctx *Ctx, r *http.Request) *Request {
	tmp := new(Request)
	tmp.ctx = ctx
	tmp.r = r
	tmp.param = tsmap.New()
	tmp.session = nil
	return tmp
}

func (this *Request) Reset(r *http.Request) {
	this.r = r
	this.query = nil
}

func (this *Request) release() {
	this.r = nil
	this.query = nil
	this.param.Release()
	this.session = nil
}

func (this *Request) Raw() *http.Request {
	return this.r
}

func (this *Request) IsMethod(method string) bool {
	return strings.EqualFold(this.r.Method, method)
}

func (this *Request) IsJSON() bool {
	return strings.Contains(this.r.Header.Get(constant.HeaderContentType), constant.MIMEApplicationJSON)
}

func (this *Request) IsXML() bool {
	return strings.Contains(this.r.Header.Get(constant.HeaderContentType), constant.MIMEApplicationXML)
}

func (this *Request) IsAjax() bool {
	return strings.EqualFold(this.r.Header.Get(constant.HeaderXRequestedWith), "XMLHttpRequest")
}

func (this *Request) AcceptJSON() bool {
	return strings.Contains(this.r.Header.Get(constant.HeaderAccept), constant.MIMEApplicationJSON)
}

func (this *Request) AcceptText() bool {
	return strings.Contains(this.r.Header.Get(constant.HeaderAccept), constant.MIMETextHTML) ||
		strings.Contains(this.r.Header.Get(constant.HeaderAccept), constant.MIMETextPlain)
}

func (this *Request) Scheme() string {
	if this.r.URL.Scheme != "" {
		return this.r.URL.Scheme
	}
	if this.r.TLS == nil {
		return "http"
	}
	return "https"
}

func (this *Request) Host() string {
	if this.r.Host == "" {
		return "localhost"
	}
	if host, _, err := net.SplitHostPort(this.r.Host); err == nil {
		return host
	}
	return this.r.Host
}

func (this *Request) Port() string {
	if this.r.Host == "" {
		return "80"
	}
	if _, port, err := net.SplitHostPort(this.r.Host); err == nil {
		return port
	}
	return "80"
}
