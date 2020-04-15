package router

import (
	"ginana-blog/internal/config"
	"ginana-blog/internal/model"
	"ginana-blog/internal/server/controller/admin"
	"ginana-blog/internal/server/controller/api"
	"ginana-blog/internal/server/controller/front"
	"ginana-blog/internal/service"
	"ginana-blog/library/mdw"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"time"
)

func InitRouter(svc service.Service, cfg *config.Config) (e *iris.Application) {
	e = NewIris(svc, cfg)
	sessManager := sessions.New(sessions.Config{
		Cookie:  "GiNana_Session",
		Expires: 24 * time.Hour,
	})
	frontParty := mvc.New(e.Party("/"))
	frontParty.Register(svc, sessManager.Start, getPagination)
	frontParty.Router.Layout("layouts/front.html")
	frontParty.Handle(new(front.CFront))
	adminParty := mvc.New(e.Party("/"))
	adminParty.Register(svc)
	adminParty.Router.Layout("layouts/admin.html")
	adminParty.Handle(new(admin.CAdmin))
	apiParty := mvc.New(e.Party("/api", mdw.Cors()).AllowMethods(iris.MethodOptions))
	apiParty.Register(svc)
	apiParty.Handle(new(api.CApi))
	return
}

func getPagination(ctx iris.Context) *model.Pagination {
	return &model.Pagination{
		Page:     ctx.URLParamInt64Default("page", 1),
		PageSize: ctx.URLParamInt64Default("pagesize", 10),
	}
}
