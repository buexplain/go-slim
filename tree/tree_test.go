package tree

import (
	"testing"
)

type Param map[string]string

func (this Param) Set(k string, v interface{}) {
	this[k] = v.(string)
}

func (this Param) Get(k string) string {
	return this[k]
}

//测试静态路由
func TestStaticUrl(t *testing.T) {
	tree := New(false)
	var param Param

	if err := tree.Add("/", 1); err != nil {
		t.Fatal("添加路由失败")
	}

	if err := tree.Add("/article", 2); err != nil {
		t.Fatal("添加路由失败")
	}

	if err := tree.Add("/article/:id", 3); err != nil {
		t.Fatal("添加路由失败")
	}

	if err := tree.Add("/content/:id.html", 4); err != nil {
		t.Fatal("添加路由失败")
	}

	param = make(Param)
	data, ok := tree.Search("/", param)
	if !ok {
		t.Fatal("静态的url / 查找失败")
	}
	if data.(int) != 1 {
		t.Fatal("静态的url / 查找失败")
	}

	data, ok = tree.Search("/article", param)
	if !ok {
		t.Fatal("静态的url /article 查找失败")
	}
	if data.(int) != 2 {
		t.Fatal("静态的url /article 查找失败")
	}

	data, ok = tree.Search("/article/100", param)
	if !ok {
		t.Fatal("静态的url /article/:id 查找失败")
	}
	if data.(int) != 3 {
		t.Fatal("静态的url /article/:id 查找失败")
	}
	if param.Get("id") != "100" {
		t.Fatal("静态的url /article/:id 查找失败")
	}

	data, ok = tree.Search("/content/10.html", param)
	if !ok {
		t.Fatal("静态的url /content/:id.html 查找失败")
	}
	if data.(int) != 4 {
		t.Fatal("静态的url /content/:id.html 查找失败")
	}
	if tmp := param.Get("id.html"); tmp != "10.html" {
		t.Fatal("静态的url /content/:id.html 查找失败 " + tmp)
	}
}

//测试动态路由
func TestDynamicUrl(t *testing.T) {
	tree := New(false)
	var param Param

	if err := tree.Add("/article/:id", 1); err != nil {
		t.Fatal("添加路由失败")
	}
	if err := tree.Add("/article/:id/status/:status", 2); err != nil {
		t.Fatal("添加路由失败")
	}
	if err := tree.Add("/backend/article/attachment/check/:m-d_5", 3); err != nil {
		t.Fatal("添加路由失败")
	}

	param = make(Param)
	data, ok := tree.Search("/article/1", param)
	if !ok {
		t.Fatal("动态参数的url查找失败")
	}
	if data.(int) != 1 {
		t.Fatal("动态参数的url查找失败", data)
	}
	if id, ok := param["id"]; !ok || id != "1" {
		t.Error("动态参数测试失败")
	}

	data, ok = tree.Search("/article/2/status/3", param)
	if !ok {
		t.Fatal("动态参数的url查找失败")
	}
	if data.(int) != 2 {
		t.Fatal("动态参数的url查找失败", data)
	}
	if id, ok := param["id"]; !ok || id != "2" {
		t.Error("动态 id 参数测试失败")
	}
	if status, ok := param["status"]; !ok || status != "3" {
		t.Error("动态 status 参数测试失败")
	}

	data, ok = tree.Search("/backend/article/attachment/check/33e2ac4290c7d44f4ea727641d7f1cff", param)
	if !ok {
		t.Fatal("动态参数的url查找失败")
	}
	if data.(int) != 3 {
		t.Fatal("动态参数的url查找失败", data)
	}
	if id, ok := param["m-d_5"]; !ok || id != "33e2ac4290c7d44f4ea727641d7f1cff" {
		t.Error("动态参数测试失败")
	}
}
