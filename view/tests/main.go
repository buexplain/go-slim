package main

import (
	"fmt"
	"github.com/buexplain/go-slim/view"
	"os"
)

type User struct {
	ID       int
	Account  string
	Password string
}

var v *view.View

func init() {
	v = view.New("./view/test/tpl", false)
}

func testEdit() {
	user := User{ID: 1, Account: "test", Password: "123456"}
	var data map[string]interface{} = map[string]interface{}{}
	data["user"] = user
	err := v.Render(os.Stdout, "edit.html", data)
	if err != nil {
		fmt.Println(err)
	}
}

func testCreate() {
	err := v.Render(os.Stdout, "create.html", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func testJS() {
	err := v.Render(os.Stdout, "layout.js.html", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	testEdit()
}
