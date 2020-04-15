package admin

import (
	"ginana-blog/internal/service"
	"github.com/kataras/iris/v12"
)

type CAdmin struct {
	Ctx iris.Context
	Svc service.Service
}
