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

func InitRouter(svc service.Service, cfg *config.Config) (e *iris.Application, err error) {

	e = newIris(cfg)

	session := sessions.New(sessions.Config{
		Cookie:  "GiNana_Session",
		Expires: 24 * time.Hour,
	})

	group := mvc.New(e.Party("/"))

	group.HandleError(func(ctx iris.Context, err error) {
		ctx.ViewData("disableRight", true)
		ctx.ViewData("error", model.PlusJson(nil, err))
		ctx.View("error/error.html")
	})

	group.Register(
		svc, session.Start,
		getPagination,
		getSiteOptions(svc, cfg),
	)

	group.Router.Layout("layouts/front.html")
	group.Handle(new(front.CFront))
	adminParty := group.Party("/admin")
	adminParty.Router.Layout("layouts/admin.html")
	adminParty.Handle(new(admin.CAdmin))

	apiParty := mvc.New(e.Party("/api", mdw.CORS([]string{"*"})).AllowMethods(iris.MethodOptions))
	apiParty.Register(svc, getPagination)
	apiParty.Handle(new(api.CApi))

	return
}
