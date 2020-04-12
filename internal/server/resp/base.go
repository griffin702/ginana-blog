package resp

import (
	"ginana-blog/library/ecode"
	"github.com/kataras/iris/v12/mvc"
)

type JSON struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func PlusJson(data interface{}, err error) *JSON {
	ec := ecode.Cause(err)
	return &JSON{
		Code:    ec.Code(),
		Message: ec.Message(),
		Data:    data,
	}
}

func PlusHtml(data map[string]interface{}, err error) mvc.Result {
	ec := ecode.Cause(err)
	data["error"] = &JSON{
		Code:    ec.Code(),
		Message: ec.Message(),
	}
	return mvc.View{
		Name: "error/error.html",
		Data: data,
	}
}
