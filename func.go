package slim

import (
	"fmt"
	"github.com/buexplain/go-slim/constant"
	"github.com/buexplain/go-slim/errors"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
)

//恐慌恢复
func defaultRecoverFunc(ctx *Ctx, a interface{}) {
	if err, ok := a.(interface{ Error() string }); ok {
		markerErr := errors.IsMarker(err)
		if markerErr == nil {
			//未知的错误，转为服务端错误，并加上栈信息
			err = fmt.Errorf("%w\n%s", err, debug.Stack())
			markerErr = errors.Mark(err, errors.ServerCode).(*errors.MrKErr)
		}
		if markerErr.Code() >= errors.ServerCode {
			defaultServerErrorFunc(ctx, markerErr)
		} else {
			defaultClientErrorFunc(ctx, markerErr)
		}
	} else {
		//未知的恐慌，转为服务端错误，并加上栈信息
		err := fmt.Errorf("%+v\n%s", a, debug.Stack())
		markerErr := errors.MarkServer(err).(*errors.MrKErr)
		defaultServerErrorFunc(ctx, markerErr)
	}
}

//服务端错误处理
func defaultServerErrorFunc(ctx *Ctx, markerErr *errors.MrKErr) {
	ctx.Response().Buffer().Reset()
	isDebug := ctx.App().Debug()
	isJSON := (!ctx.Request().AcceptText() || (ctx.Route() != nil && ctx.Route().HasLabel("json")))
	var responseErr error
	if isJSON {
		//返回json
		if isDebug {
			//返回具体错误
			responseErr = ctx.Response().Error(markerErr.Code(), markerErr.Error(), http.StatusOK)
		} else {
			//屏蔽错误
			responseErr = ctx.Response().Error(markerErr.Code(), http.StatusText(http.StatusInternalServerError), http.StatusOK)
		}
	} else {
		//返回文本
		ctx.Response().Header().Set(constant.HeaderXContentTypeOptions, "nosniff")
		if isDebug {
			responseErr = ctx.Response().Abort(
				http.StatusInternalServerError,
				strings.ReplaceAll(strings.ReplaceAll(markerErr.Error(), "\n", "<br>"), "\t", "&nbsp;&nbsp;&nbsp;&nbsp;"))
		} else {
			responseErr = ctx.Response().Abort(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
	}
	if !isDebug {
		//生产环境，记录错误日志
		log.Println(fmt.Sprintf(
			"%s%s %s%s %s%d%s %s",
			"[",
			ctx.Request().Raw().Method,
			ctx.Request().Raw().URL.String(),
			"]",
			"[",
			markerErr.Code(),
			"]",
			markerErr.Error(),
		))
	}
	//响应失败，记录日志
	if responseErr != nil {
		log.Println(fmt.Sprintf(
			"%s%s %s%s %s%d%s %s",
			"[",
			ctx.Request().Raw().Method,
			ctx.Request().Raw().URL.String(),
			"]",
			"[",
			markerErr.Code(),
			"]",
			responseErr.Error(),
		))
	}
}

//客户端错误处理
func defaultClientErrorFunc(ctx *Ctx, markerErr *errors.MrKErr) {
	ctx.Response().Buffer().Reset()
	isJSON := (!ctx.Request().AcceptText() || (ctx.Route() != nil && ctx.Route().HasLabel("json")))
	var responseErr error
	if isJSON {
		//返回json
		responseErr = ctx.Response().Error(markerErr.Code(), markerErr.Error(), http.StatusOK)
	} else {
		//返回文本
		ctx.Response().Header().Set(constant.HeaderXContentTypeOptions, "nosniff")
		responseErr = ctx.Response().Abort(http.StatusBadRequest, markerErr.Error())
	}
	//响应失败，记录日志
	if responseErr != nil {
		log.Println(fmt.Sprintf(
			"%s%s %s%s %s%d%s %s",
			"[",
			ctx.Request().Raw().Method,
			ctx.Request().Raw().URL.String(),
			"]",
			"[",
			markerErr.Code(),
			"]",
			responseErr.Error(),
		))
	}
}

//错误处理
func defaultErrorFunc(ctx *Ctx, err error) {
	if err == nil {
		return
	}
	markerErr := errors.IsMarker(err)
	if markerErr == nil {
		//未知错误，转为服务端错误
		markerErr = errors.Mark(err, errors.ServerCode).(*errors.MrKErr)
	}
	if markerErr.Code() >= errors.ServerCode {
		defaultServerErrorFunc(ctx, markerErr)
	} else {
		defaultClientErrorFunc(ctx, markerErr)
	}
}

//默认路由错误处理
func defaultRoute(ctx *Ctx, w *Response, r *Request) error {
	ctx.Response().Buffer().Reset()
	isJSON := (!ctx.Request().AcceptText() || (ctx.Route() != nil && ctx.Route().HasLabel("json")))
	if isJSON {
		//返回json
		return ctx.Response().Error(errors.ClientCode, "404 route not found", http.StatusOK)
	} else {
		//返回文本
		ctx.Response().Header().Set(constant.HeaderXContentTypeOptions, "nosniff")
		return w.HTML(http.StatusNotFound, "404 route not found")
	}
}
