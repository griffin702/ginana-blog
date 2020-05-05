package server

import (
	"ginana-blog/internal/config"
	"ginana-blog/internal/controller/admin"
	"ginana-blog/internal/controller/api"
	"ginana-blog/internal/controller/front"
	"ginana-blog/internal/model"
	"ginana-blog/internal/service"
	"github.com/griffin702/ginana/library/mdw"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"time"
)

func InitRouter(svc service.Service, cfg *config.Config, hm service.HelperMap, valid model.ValidatorHandler) (e *iris.Application, err error) {

	e = newIris(cfg)

	session := sessions.New(sessions.Config{
		Cookie:  "GiNana_Session",
		Expires: time.Duration(cfg.SessionAndCookieExpire),
	})

	objects := []interface{}{
		svc, session.Start, hm, valid, jsonPlus,
		getClientIP, getPagination, getTools, getConfigs,
		getSiteOptions(svc, cfg),
	}

	group := mvc.New(e.Party("/"))

	group.HandleError(func(ctx iris.Context, err error) {
		ctx.ViewData("disableRight", true)
		ctx.ViewData("error", jsonPlus(ctx)(nil, err))
		ctx.View("error/error.html")
	})

	group.Register(objects...)

	group.Router.Layout("layouts/front.html")
	group.Handle(new(front.CFront))
	adminParty := group.Party("/admin")
	adminParty.Router.Layout("layouts/admin.html")
	adminParty.Handle(new(admin.CAdmin))

	apiParty := mvc.New(e.Party("/api", mdw.CORS([]string{"*"})).AllowMethods(iris.MethodOptions))
	apiParty.Register(objects...)
	apiParty.Handle(new(api.CApi))

	return
}
