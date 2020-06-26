package slim

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/buexplain/go-slim/constant"
	"github.com/buexplain/go-slim/errors"
	"github.com/buexplain/go-slim/tsmap"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Response struct {
	ctx        *Ctx
	w          http.ResponseWriter
	store      *tsmap.TSMap
	statusCode int
	buffer     *bytes.Buffer
}

func NewResponse(ctx *Ctx, w http.ResponseWriter) *Response {
	tmp := new(Response)
	tmp.ctx = ctx
	tmp.w = w
	tmp.store = tsmap.New()
	tmp.statusCode = 0
	tmp.buffer = new(bytes.Buffer)
	return tmp
}

func (this *Response) Buffer() *bytes.Buffer {
	return this.buffer
}

func (this *Response) Reset(w http.ResponseWriter) {
	this.w = w
}

func (this *Response) release() {
	this.w = nil
	this.statusCode = 0
	this.store.Release()
	this.buffer.Reset()
}

func (this *Response) send() bool {
	if this.statusCode != 0 {
		//先写header
		if this.ctx.r.session != nil {
			if err := this.ctx.r.session.Save(this.ctx.r.r, this); err != nil {
				//将session设置为nil，避免死循环
				this.ctx.r.session = nil
				//抛出错误
				this.ctx.Throw(fmt.Errorf("save session error: %w", err))
				return false
			}
		}
		//再写code
		this.w.WriteHeader(this.statusCode)
		//最后写body
		if this.buffer.Len() > 0 {
			_, _ = this.buffer.WriteTo(this.w)
		}
	}
	return true
}

func (this *Response) Raw() http.ResponseWriter {
	return this.w
}

func (this *Response) Header() http.Header {
	return this.w.Header()
}

func (this *Response) Write(b []byte) (int, error) {
	i, err := this.buffer.Write(b)
	if err != nil {
		return i, err
	}
	return i, nil
}

func (this *Response) WriteHeader(statusCode int) {
	this.statusCode = statusCode
}

func (this *Response) StatusCode() int {
	return this.statusCode
}

func (this *Response) Redirect(statusCode int, url string) error {
	http.Redirect(this, this.ctx.r.r, url, statusCode)
	return nil
}

func (this *Response) RedirectBack() error {
	tmp := this.ctx.r.r.Header.Get("Referer")
	if tmp == "" {
		tmp = "/"
	}
	http.Redirect(this, this.ctx.r.r, tmp, http.StatusFound)
	return nil
}

func (this *Response) Plain(statusCode int, text string) error {
	this.w.Header().Set(constant.HeaderContentType, constant.MIMETextPlainCharsetUTF8)
	this.WriteHeader(statusCode)
	_, err := this.Write([]byte(text))
	return err
}

func (this *Response) HTML(statusCode int, text string) error {
	this.w.Header().Set(constant.HeaderContentType, constant.MIMETextHTMLCharsetUTF8)
	this.WriteHeader(statusCode)
	_, err := this.Write([]byte(text))
	return err
}

func (this *Response) Render(wr io.Writer, data interface{}, tpl string) error {
	return this.ctx.app.view.Render(wr, tpl, data)
}

func (this *Response) Store() *tsmap.TSMap {
	return this.store
}

func (this *Response) Assign(key string, a interface{}) *Response {
	if key == "" {
		panic(errors.New("assign key not allow empty"))
	}
	this.store.Set(key, a)
	return this
}

func (this *Response) Abort(statusCode int, message ...interface{}) error {
	tpl := "errors/" + strconv.Itoa(statusCode) + ".html"
	if len(message) > 0 {
		this.store.Set("message", message[0])
	} else {
		this.store.Set("message", http.StatusText(statusCode))
	}
	return this.View(statusCode, tpl)
}

func (this *Response) Jump(url string, message interface{}, wait ...int) error {
	this.store.Set("url", url)
	if len(wait) > 0 {
		this.store.Set("wait", wait[0])
	} else {
		this.store.Set("wait", 5)
	}
	return this.Abort(http.StatusFound, message)
}

func (this *Response) JumpBack(message interface{}, wait ...int) error {
	tmp := this.ctx.r.r.Header.Get("Referer")
	if tmp == "" {
		tmp = "/"
	}
	return this.Jump(tmp, message, wait...)
}

func (this *Response) View(statusCode int, tpl string) error {
	if len(tpl) == 0 {
		panic("view template not allow empty")
	}
	buff := &bytes.Buffer{}
	err := this.ctx.app.view.Render(buff, tpl, this.store.Pop())
	if err != nil {
		return err
	}
	if buff.Len() == 0 {
		return fmt.Errorf("template: %s undefined", tpl)
	}
	this.w.Header().Set(constant.HeaderContentType, constant.MIMETextHTMLCharsetUTF8)
	this.WriteHeader(statusCode)
	_, err = this.Write(buff.Bytes())
	return err
}

//返回一个json
func (this *Response) JSON(statusCode int) error {
	var content []byte
	var err error
	content, err = json.Marshal(this.store.Pop())
	if err != nil {
		return err
	}
	this.w.Header().Set(constant.HeaderContentType, constant.MIMEApplicationJSONCharsetUTF8)
	this.WriteHeader(statusCode)
	_, err = this.Write(content)
	return err
}

//返回成功的json
func (this *Response) Success(data ...interface{}) error {
	if l := len(data); l == 0 {
		if !this.store.Has("data") {
			this.Assign("data", "")
		}
	} else {
		if l == 1 {
			this.Assign("data", data[0])
		} else {
			this.Assign("data", data)
		}
	}
	if !this.store.Has("code") {
		this.Assign("code", 0)
	}
	if !this.store.Has("message") {
		this.Assign("message", "success")
	}
	return this.JSON(http.StatusOK)
}

//返回错误的json
func (this *Response) Error(code int, message string, statusCode ...int) error {
	if !this.store.Has("data") {
		this.Assign("data", "")
	}
	this.Assign("code", code).Assign("message", message)
	if len(statusCode) == 0 {
		//一般是客户端错误，所以默认为200
		return this.JSON(http.StatusOK)
	}
	return this.JSON(statusCode[0])
}

func (this *Response) JSONP(statusCode int, callback string) error {
	var content []byte
	var err error
	content, err = json.Marshal(this.store.Pop())
	if err != nil {
		return err
	}
	callback = template.JSEscapeString(callback)
	this.w.Header().Set(constant.HeaderContentType, constant.MIMEApplicationJavaScriptCharsetUTF8)
	this.WriteHeader(statusCode)
	buff := bytes.NewBufferString(" if(window.")
	buff.WriteString(callback)
	buff.WriteByte(')')
	buff.WriteString(callback)
	buff.WriteByte('(')
	buff.Write(content)
	buff.WriteString(");")
	_, err = this.Write(buff.Bytes())
	return err
}

func (this *Response) XML(statusCode int, x interface{}) error {
	var content []byte
	var err error
	content, err = xml.Marshal(x)
	if err != nil {
		return err
	}
	this.w.Header().Set(constant.HeaderContentType, constant.MIMEApplicationXMLCharsetUTF8)
	this.WriteHeader(statusCode)
	_, err = this.Write(content)
	return err
}

func (this *Response) Cookie(name string, value string, argv ...interface{}) *Response {
	cookie := http.Cookie{}
	cookie.Name = name
	cookie.Value = url.QueryEscape(value)

	l := len(argv)

	cookie.MaxAge = 3600
	if l > 0 {
		if v, ok := argv[0].(int); ok {
			cookie.MaxAge = v
		}
	}

	cookie.Path = "/"
	if l > 2 {
		if v, ok := argv[2].(string); ok && len(v) > 0 {
			cookie.Path = v
		}
	}

	if l > 2 {
		if v, ok := argv[2].(string); ok && len(v) > 0 {
			cookie.Domain = v
		}
	}

	cookie.Secure = false
	if l > 3 {
		if v, ok := argv[3].(bool); ok {
			cookie.Secure = v
		}
	}

	cookie.HttpOnly = true
	if l > 4 {
		if v, ok := argv[4].(bool); ok {
			cookie.HttpOnly = v
		}
	}

	this.Header().Add(constant.HeaderSetCookie, cookie.String())

	return this
}

func (this *Response) File(file string) error {
	if strings.Contains(this.ctx.r.r.URL.Path, "..") {
		for _, ent := range strings.FieldsFunc(this.ctx.r.r.URL.Path, func(r rune) bool { return r == '/' || r == '\\' }) {
			if ent == ".." {
				return errors.MarkClient(errors.New("invalid URL path"))
			}
		}
	}

	dir, name := filepath.Split(file)
	fs := http.Dir(dir)
	f, err := fs.Open(name)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.MarkClient(errors.New("404 file not found"))
		}
		if os.IsPermission(err) {
			return err
		}
		return errors.MarkClient(fmt.Errorf("invalid URL path: %w", err))
	}

	defer func() {
		_ = f.Close()
	}()

	var fi os.FileInfo
	fi, err = f.Stat()
	if err != nil {
		if os.IsNotExist(err) {
			return errors.MarkClient(errors.New("404 file not found"))
		}
		if os.IsPermission(err) {
			return err
		}
		return errors.MarkClient(fmt.Errorf("invalid URL path: %w", err))
	}

	if fi.IsDir() {
		return errors.MarkClient(errors.New("403 Forbidden"))
	}

	http.ServeContent(this, this.ctx.r.r, fi.Name(), fi.ModTime(), f)
	return nil
}

func (this *Response) Download(file string, name ...string) error {
	var fName string
	if len(name) > 0 && name[0] != "" {
		fName = name[0]
	} else {
		_, fName = filepath.Split(file)
	}
	this.Header().Add("Content-Disposition", "attachment; filename="+url.QueryEscape(fName))
	this.Header().Add("Content-Description", "File Transfer")
	this.Header().Add(constant.HeaderContentType, constant.MIMEOctetStream)
	this.Header().Add("Content-Transfer-Encoding", "binary")
	this.Header().Add("Expires", "Mon, 26 Jul 1997 05:00:00 GMT")
	this.Header().Add("Cache-Control", "cache, must-revalidate")
	this.Header().Add("Pragma", "public")
	return this.File(file)
}
