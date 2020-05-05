package api

import (
	"ginana-blog/internal/model"
	"ginana-blog/internal/service"
	"github.com/griffin702/service/tools"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

type CApi struct {
	Ctx         iris.Context
	Session     *sessions.Session
	Svc         service.Service
	Pager       *model.Pager
	GetClientIP model.GetClientIP
	Hm          service.HelperMap
	Valid       model.Validator
	Tool        *tools.Tool
}
