package main

import (
	"bytes"
	"fmt"
	"github.com/buexplain/go-slim"
	"github.com/buexplain/go-slim/constant"
	"github.com/buexplain/go-slim/errors"
	"io"
	"log"
	"net/http"
)

//http应用
var app *slim.App

var addr string

func init() {
	//初始化app
	app = slim.New(true)
	addr = "127.0.0.1:1991"
}

func main() {
	//支持全局中间件
	app.Use(func(ctx *slim.Ctx, w *slim.Response, r *slim.Request) {
		log.Println("进入全局中间件")
		ctx.Next()
	})

	app.Mux().Any("error", func(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
		act := r.Query("act", "abort")
		switch act {
		case "abort":
			return w.Abort(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		case "panicStr":
			panic("test panic string")
		case "panicRawErr":
			panic(errors.New("test panic raw err"))
		case "panicMrKErrServer":
			panic(errors.MarkServer(errors.New("test panic mark server err")))
		case "panicMrKErrClient":
			panic(errors.MarkClient(errors.New("test panic mark client err")))
		case "jump":
			return w.Jump("/", "test jump to index", 5)
		case "returnRawErr":
			return errors.New("test return raw err")
		case "returnMrKErrServer":
			return errors.MarkServer(errors.New("test return mark server err"))
		case "returnMrKErrClient":
			return errors.MarkClient(errors.New("test return mark client err"))
		}
		return w.Abort(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	})

	//支持响应缓冲
	app.Mux().Any("/", func(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
		buf := &bytes.Buffer{}
		_, _ = io.WriteString(buf, "<h3>点击如下按钮进行测试，注意观察命令行输出信息：</h3><hr>")
		_, _ = io.WriteString(buf, "<a target='_blank' href='/article/1'>全局、组、xml响应</a><br>")
		_, _ = io.WriteString(buf, "<a target='_blank' href='/user/1'>全局、组、路由中间件、json响应</a><br>")
		_, _ = io.WriteString(buf, "<a target='_blank' href='/error?act=abort'>abort</a><br>")
		_, _ = io.WriteString(buf, "<a target='_blank' href='/error?act=panicStr'>panic str</a><br>")
		_, _ = io.WriteString(buf, "<a target='_blank' href='/error?act=panicRawErr'>panic raw err</a><br>")
		_, _ = io.WriteString(buf, "<a target='_blank' href='/error?act=panicMrKErrServer'>panic server mrKErr</a><br>")
		_, _ = io.WriteString(buf, "<a target='_blank' href='/error?act=panicMrKErrClient'>panic client mrKErr</a><br>")
		_, _ = io.WriteString(buf, "<a target='_blank' href='/error?act=returnRawErr'>return raw err</a><br>")
		_, _ = io.WriteString(buf, "<a target='_blank' href='/error?act=returnMrKErrServer'>return server mrKErr</a><br>")
		_, _ = io.WriteString(buf, "<a target='_blank' href='/error?act=returnMrKErrClient'>return client mrKErr</a><br>")
		_, _ = io.WriteString(buf, "<a target='_blank' href='/error?act=jump'>jump</a><br>")
		//先写数据
		_, err := io.Copy(w, buf)
		//再写header头
		w.WriteHeader(http.StatusOK)
		w.Header().Set(constant.HeaderContentType, constant.MIMETextHTMLCharsetUTF8)
		//返回错误
		return err
	})

	//支持组路由
	app.Mux().Group("", func() {
		//文章详情
		app.Mux().Get("/article/:id", func(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
			type Result struct {
				Data string
				Code int
				Msg  string
			}
			return w.XML(http.StatusOK, Result{Data: r.Input("id", "0"), Code: 0, Msg: "success"})
		})
		//用户详情
		app.Mux().Get("/user/:id", func(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
			return w.Assign("code", 0).Assign("data", r.Input("id", "0")).Assign("msg", "success").JSON(http.StatusOK)
		}).Use(func(ctx *slim.Ctx, w *slim.Response, r *slim.Request) { //支持路由中间件
			log.Println("进入路由中间件")
			ctx.Next()
		})
	}).Use(func(ctx *slim.Ctx, w *slim.Response, r *slim.Request) {
		log.Println("进入路由组中间件")
		ctx.Next()
	}).Regexp("id", `\d`) //用正则约束id

	//打印路由列表
	fmt.Println("http://" + addr + "/")
	fmt.Println(app.Mux().DumpRouteMap("main."))

	//启动服务器
	log.Fatalln(app.Run(addr))
}
