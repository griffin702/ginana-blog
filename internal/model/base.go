package model

import (
	"ginana-blog/library/ecode"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"net/url"
	"strings"
)

// BlogGin hello BlogGin.
type GiNana struct {
	Hello string
}

type JSON struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Rawurlencode(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

func PlusJson(data interface{}, err error) *JSON {
	ec := ecode.Cause(err)
	return &JSON{
		Code:    ec.Code(),
		Message: ec.Message(),
		Data:    data,
	}
}

func PlusHtmlErr(ctx iris.Context, err error) mvc.Result {
	ec := ecode.Cause(err)
	ctx.ViewData("error", &JSON{
		Code:    ec.Code(),
		Message: ec.Message(),
	})
	return mvc.View{
		Name: "error/error.html",
	}
}
