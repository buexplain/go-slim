package view

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template/parse"
)

type View struct {
	path           string
	leftDelimiter  string
	rightDelimiter string
	funcMap        template.FuncMap
	cache          map[string]*template.Template
	isCache        bool
	l              *sync.RWMutex
}

func New(path string, isCache bool) *View {
	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	path = filepath.ToSlash(path)
	tmp := new(View)
	tmp.path = path
	tmp.leftDelimiter = "{{"
	tmp.rightDelimiter = "}}"
	tmp.funcMap = template.FuncMap{}
	tmp.isCache = isCache
	tmp.cache = make(map[string]*template.Template)
	tmp.l = new(sync.RWMutex)
	tmp.AddFunc("HTML", HTML)
	return tmp
}

func (this *View) AddFunc(name string, f interface{}) *View {
	if name == "" {
		panic(fmt.Errorf("template func name not allow empty"))
	}
	if _, ok := this.funcMap[name]; !ok {
		this.funcMap[name] = f
	}
	return this
}

func (this *View) SetCache(isCache bool) *View {
	this.l.Lock()
	defer this.l.Unlock()
	if this.isCache != isCache {
		this.isCache = isCache
		this.cache = make(map[string]*template.Template)
	}
	return this
}

func (this *View) SetPath(path string) *View {
	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	this.path = path
	return this
}

//读取模板文件
func (this *View) getFile(tpl string) ([]byte, error) {
	if len(tpl) == 0 {
		return nil, fmt.Errorf("template name not allow empty")
	}
	if _, err := os.Stat(tpl); err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(tpl)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//获取一个模板
func (this *View) getTemplate(tpl string) (*template.Template, error) {
	b, err := this.getFile(tpl)
	if err != nil {
		return nil, err
	}
	t := template.New(tpl).Delims(this.leftDelimiter, this.rightDelimiter).Funcs(this.funcMap)
	t, err = t.Parse(string(b))
	if err != nil {
		return nil, err
	}
	return t, err
}

//解析继承的模板
func (this *View) parseExtend(tpl string, depth int, preTpl string) ([]string, error) {
	if depth > 10 {
		if preTpl == "" {
			return nil, fmt.Errorf("cycle include is not allowed for template: %s", tpl)
		} else {
			return nil, fmt.Errorf("cycle include is not allowed for template: %s <=> %s", preTpl, tpl)
		}
	}
	t, err := this.getTemplate(tpl)
	if err != nil {
		return nil, err
	}
	tArr := t.Templates()
	extendArr := make([]string, 0)
	for _, t := range tArr {
		if t.Name() == "extend" {
			extend := filepath.ToSlash(filepath.Join(this.path, strings.TrimSpace(t.Tree.Root.String())))
			tmp, err := this.parseExtend(extend, depth+1, tpl)
			if err != nil {
				return nil, err
			}
			if len(tmp) > 0 {
				extendArr = append(extendArr, tmp...)
			}
			extendArr = append(extendArr, extend)
			break
		}
	}
	return extendArr, nil
}

//解析引入的文件
func (this *View) parseInclude(tpl string, depth int, preTpl string) ([]string, error) {
	if depth > 10 {
		if preTpl == "" {
			return nil, fmt.Errorf("cycle include is not allowed for template: %s", tpl)
		} else {
			return nil, fmt.Errorf("cycle include is not allowed for template: %s <=> %s", preTpl, tpl)
		}
	}
	t, err := this.getTemplate(tpl)
	if err != nil {
		return nil, err
	}
	tArr := t.Templates()
	includeArr := make([]string, 0)
	for _, t := range tArr {
		if t.Tree == nil || t.Tree.Root == nil || t.Tree.Root.Nodes == nil {
			continue
		}
		for _, v := range t.Tree.Root.Nodes {
			if v.Type() != parse.NodeTemplate {
				continue
			}
			node, ok := v.(*parse.TemplateNode)
			if !ok {
				continue
			}
			include := filepath.ToSlash(filepath.Join(this.path, node.Name))
			//跳过不存在的模板
			if _, err := os.Stat(include); err != nil {
				continue
			}
			tmp, err := this.parseInclude(include, depth+1, tpl)
			if err != nil {
				return nil, err
			}
			includeArr = append(includeArr, include)
			if len(tmp) > 0 {
				includeArr = append(includeArr, tmp...)
			}
		}
	}
	return includeArr, nil
}

//[]string去重
func (this View) uniqueStrSlice(strSlice []string) []string {
	result := make([]string, 0, len(strSlice))
	unique := make(map[string]bool)
	for _, v := range strSlice {
		if _, ok := unique[v]; ok {
			continue
		}
		unique[v] = true
		result = append(result, v)
	}
	return result
}

//解析模板的依赖
func (this *View) parseDependent(tpl string) ([]string, error) {
	tpl = filepath.ToSlash(filepath.Join(this.path, tpl))
	extendArr, err := this.parseExtend(tpl, 0, "")
	if err != nil {
		return nil, err
	}
	extendArr = append(extendArr, tpl)
	dependent := make([]string, 0)
	for _, extend := range extendArr {
		include, err := this.parseInclude(extend, 0, "")
		if err != nil {
			return nil, err
		}
		dependent = append(dependent, extend)
		dependent = append(dependent, include...)
	}
	dependent = this.uniqueStrSlice(dependent)
	return dependent, nil
}

//解析模板
func (this *View) parseTemplate(tpl string) (*template.Template, error) {
	dependent, err := this.parseDependent(tpl)
	if err != nil {
		return nil, err
	}
	var t *template.Template
	for _, filename := range dependent {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		s := string(b)
		name := filename[len(this.path)+1:]
		var tmpl *template.Template
		if t == nil {
			t = template.New(name).Delims(this.leftDelimiter, this.rightDelimiter).Funcs(this.funcMap)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name).Delims(this.leftDelimiter, this.rightDelimiter).Funcs(this.funcMap)
		}
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func (this *View) parseTemplateFromCache(tpl string) (*template.Template, error) {
	this.l.RLock()
	if t, ok := this.cache[tpl]; ok {
		this.l.RUnlock()
		return t, nil
	}
	this.l.RUnlock()

	this.l.Lock()
	defer this.l.Unlock()

	if t, ok := this.cache[tpl]; ok {
		return t, nil
	}

	t, err := this.parseTemplate(tpl)
	if err == nil {
		this.cache[tpl] = t
	}
	return t, err
}

//渲染模板
func (this *View) Render(wr io.Writer, tpl string, data interface{}) error {
	var t *template.Template
	var err error
	if this.isCache {
		t, err = this.parseTemplateFromCache(tpl)
	} else {
		t, err = this.parseTemplate(tpl)
	}
	if err != nil {
		return fmt.Errorf("view render error: %w", err)
	}
	return t.Execute(wr, data)
}
