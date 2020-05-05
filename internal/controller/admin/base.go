package admin

import (
	"ginana-blog/internal/config"
	"ginana-blog/internal/model"
	"ginana-blog/internal/service"
	"github.com/griffin702/service/tools"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

type CAdmin struct {
	Ctx         iris.Context
	Session     *sessions.Session
	Svc         service.Service
	JsonPlus    model.JsonPlus
	Pager       *model.Pager
	GetClientIP model.GetClientIP
	GetOption   model.GetOption
	Links       *model.Links
	Hm          service.HelperMap
	Valid       model.Validator
	Tool        *tools.Tool
	Config      *config.Config
}
