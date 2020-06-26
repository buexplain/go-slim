package slim

import (
	"fmt"
	"github.com/buexplain/go-slim/constant"
	"github.com/buexplain/go-slim/errors"
	"github.com/buexplain/go-slim/upload"
	"github.com/gorilla/schema"
	"mime"
	"net/http"
	"strconv"
	"strings"
)

func (this *Request) ParseForm(maxMemory ...int64) error {
	if this.r.PostForm == nil && (this.r.Method == http.MethodPost || this.r.Method == http.MethodPut || this.r.Method == http.MethodPatch) {
		var err error
		ct := this.r.Header.Get(constant.HeaderContentType)
		if ct == "" {
			ct = constant.MIMEOctetStream
		}
		ct, _, err = mime.ParseMediaType(ct)
		if err != nil {
			return err
		}
		if ct == constant.MIMEMultipartForm {
			this.r.Body = http.MaxBytesReader(this.ctx.w.Raw(), this.r.Body, this.ctx.app.bodyMaxBytes)
			if len(maxMemory) != 0 && maxMemory[0] > 0 {
				err = this.r.ParseMultipartForm(maxMemory[0])
			}else {
				err = this.r.ParseMultipartForm(this.ctx.app.formMaxMemory)
			}
		} else {
			err = this.r.ParseForm()
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Request) hasForm(key string) (string, bool) {
	err := this.ParseForm()
	if err != nil {
		return key, false
	}
	if this.r.PostForm == nil {
		return key, false
	}
	_, ok := this.r.PostForm[key]
	if !ok {
		//判断是否存在带中括号类型的
		l := len(key)
		if l <= 2 || key[l-2:] != "[]" {
			tmp := key + "[]"
			_, ok = this.r.PostForm[tmp]
			if ok {
				//存在则将key重置
				key = tmp
			}
		}
	}
	return key, ok
}

func (this *Request) HasForm(key string) bool {
	_, ok := this.hasForm(key)
	return ok
}

func (this *Request) Form(key string, def ...string) string {
	var ok bool
	key, ok = this.hasForm(key)
	if ok {
		data := this.r.PostForm[key]
		if len(data) > 0 {
			return data[0]
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

func (this *Request) FormBool(key string, def ...bool) bool {
	var ok bool
	key, ok = this.hasForm(key)
	if ok {
		data := this.r.PostForm[key]
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

func (this *Request) FormInt(key string, def ...int) int {
	var ok bool
	key, ok = this.hasForm(key)
	if ok {
		data := this.r.PostForm[key]
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

func (this *Request) FormPositiveInt(key string, def ...int) int {
	var ok bool
	key, ok = this.hasForm(key)
	if ok {
		data := this.r.PostForm[key]
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

func (this *Request) FormFloat64(key string, def ...float64) float64 {
	var ok bool
	key, ok = this.hasForm(key)
	if ok {
		data := this.r.PostForm[key]
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

func (this *Request) FormSlice(key string, def ...[]string) []string {
	var ok bool
	key, ok = this.hasForm(key)
	if ok {
		result := this.r.PostForm[key]
		if len(result) > 0 {
			return result
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

func (this *Request) FormSliceBool(key string, def ...[]bool) []bool {
	data := this.FormSlice(key)
	if data != nil {
		result := make([]bool, 0, len(data))
		for _, v := range data {
			if r, err := strconv.ParseBool(v); err == nil {
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

func (this *Request) FormSliceInt(key string, def ...[]int) []int {
	data := this.FormSlice(key)
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

func (this *Request) FormSlicePositiveInt(key string, def ...[]int) []int {
	data := this.FormSlice(key)
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

func (this *Request) FormSliceFloat64(key string, def ...[]float64) []float64 {
	data := this.FormSlice(key)
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

func (this *Request) FormSliceByComma(key string, def ...[]string) []string {
	var ok bool
	key, ok = this.hasForm(key)
	if ok {
		data := this.r.PostForm[key]
		if len(data) > 0 && data[0] != "" {
			//兼容中英文逗号
			tmp := strings.Split(strings.ReplaceAll(data[0], "，", ","), ",")
			result := make([]string, 0, len(tmp))
			for _, v := range tmp {
				v = strings.TrimSpace(v)
				if v != "" {
					//跳过空字符串
					result = append(result, v)
				}
			}
			if len(result) > 0 {
				return result
			}
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

func (this *Request) FormSliceBoolByComma(key string, def ...[]bool) []bool {
	data := this.FormSliceByComma(key)
	if data != nil {
		result := make([]bool, 0, len(data))
		for _, v := range data {
			if r, err := strconv.ParseBool(v); err == nil {
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

func (this *Request) FormSliceIntByComma(key string, def ...[]int) []int {
	data := this.FormSliceByComma(key)
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

func (this *Request) FormSlicePositiveIntByComma(key string, def ...[]int) []int {
	data := this.FormSliceByComma(key)
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

func (this *Request) FormSliceFloat64ByComma(key string, def ...[]float64) []float64 {
	data := this.FormSliceByComma(key)
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
func (this *Request) FormMap(key string, def ...map[string]string) map[string]string {
	if err := this.ParseForm(); err != nil {
		if len(def) > 0 {
			return def[0]
		}
		return nil
	}
	key_len := len(key)
	var data map[string]string
	var k_l int
	for k, v := range this.r.PostForm {
		k_l = len(k)
		if len(v) > 0 && k_l-2 > key_len && k[:key_len] == key && k[key_len:key_len+1] == "[" && k[k_l-1:] == "]" {
			if data == nil {
				data = map[string]string{}
			}
			data[k[key_len+1:k_l-1]] = v[0]
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

func (this *Request) FormToStruct(a interface{}) error {
	if err := this.ParseForm(); err != nil {
		return err
	}
	if this.r.PostForm == nil {
		return errors.MarkClient(errors.New("missing form body"))
	}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	if err := decoder.Decode(a, this.r.PostForm); err != nil {
		return errors.MarkClient(fmt.Errorf("parse form to struct error: %s", err))
	} else {
		return nil
	}
}

func (this *Request) hasFile(key string) (string, bool) {
	if err := this.ParseForm(); err != nil {
		return key, false
	}
	//如果是非post请求，则可能有nil的情况
	if this.r.MultipartForm == nil {
		return key, false
	}
	_, ok := this.r.MultipartForm.File[key]
	if !ok {
		//判断是否存在带中括号类型的
		l := len(key)
		if l <= 2 || key[l-2:] != "[]" {
			tmp := key + "[]"
			_, ok = this.r.MultipartForm.File[tmp]
			if ok {
				//存在则将key重置
				key = tmp
			}
		}
	}
	return key, ok
}

func (this *Request) HasFile(key string) bool {
	_, ok := this.hasFile(key)
	return ok
}

func (this *Request) File(key string) (*upload.Upload, error) {
	var ok bool
	key, ok = this.hasFile(key)
	if !ok {
		return nil, errors.MarkClient(http.ErrMissingFile)
	}
	fhs := this.r.MultipartForm.File[key]
	if len(fhs) == 0 {
		return nil, errors.MarkClient(http.ErrMissingFile)
	}
	f, err := fhs[0].Open()
	if err != nil {
		return nil, err
	}
	return upload.New(f, fhs[0]), nil
}

func (this *Request) FileSlice(key string) (upload.Uploads, error) {
	var ok bool
	key, ok = this.hasFile(key)
	if !ok {
		return nil, errors.MarkClient(http.ErrMissingFile)
	}
	fhs := this.r.MultipartForm.File[key]
	if len(fhs) == 0 {
		return nil, errors.MarkClient(http.ErrMissingFile)
	}
	uploads := upload.Uploads{}
	for _, fh := range fhs {
		f, openErr := fh.Open()
		if openErr != nil {
			//关闭已经打开的文件
			if closeErr := uploads.Close(); closeErr != nil {
				return nil, fmt.Errorf("openErr: %w closeErr: %w", openErr, closeErr)
			}
			return nil, openErr
		}
		uploads = append(uploads, upload.New(f, fh))
	}
	return uploads, nil
}

//获取符合此类 (.+?)\\[(.+?)\\] 正则的参数
func (this *Request) FileMap(key string) (map[string]*upload.Upload, error) {
	//解析表单
	if err := this.ParseForm(); err != nil {
		return nil, err
	}
	//如果是非post请求，则可能有nil的情况
	if this.r.MultipartForm == nil {
		return nil, errors.MarkClient(http.ErrMissingFile)
	}
	key_l := len(key)
	var data map[string]*upload.Upload
	var k_l int
	for k, fhs := range this.r.MultipartForm.File {
		k_l = len(k)
		if len(fhs) > 0 && k_l-2 > key_l && k[:key_l] == key && k[key_l:key_l+1] == "[" && k[k_l-1:] == "]" {
			if data == nil {
				data = make(map[string]*upload.Upload)
			}
			f, err := fhs[0].Open()
			if err != nil {
				//关闭已经打开的文件
				for _, up := range data {
					_ = up.Close()
				}
				return nil, err
			}
			data[k[key_l+1:k_l-1]] = upload.New(f, fhs[0])
		}
	}
	if data != nil {
		return data, nil
	}
	return nil, errors.MarkClient(http.ErrMissingFile)
}
