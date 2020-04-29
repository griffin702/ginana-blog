package router

import (
	"fmt"
	"ginana-blog/internal/config"
	"ginana-blog/internal/model"
	"ginana-blog/internal/server/controller/admin"
	"ginana-blog/internal/server/controller/api"
	"ginana-blog/internal/server/controller/front"
	"ginana-blog/internal/server/resp"
	"ginana-blog/internal/service"
	"ginana-blog/library/mdw"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"time"
)

func InitRouter(svc service.Service, cfg *config.Config) (e *iris.Application, err error) {
	e = NewIris(cfg)

	//e.Use(func(ctx iris.Context) {
	//	ctx.Gzip(cfg.EnableGzip)
	//	ctx.Next()
	//})

	session := sessions.New(sessions.Config{
		Cookie:  "GiNana_Session",
		Expires: 24 * time.Hour,
	})

	group := mvc.New(e.Party("/"))

	group.HandleError(func(ctx iris.Context, err error) {
		ctx.ViewData("disableRight", true)
		ctx.ViewData("error", resp.PlusJson(nil, err))
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

	apiParty := mvc.New(e.Party("/api", mdw.Cors([]string{"*"})).AllowMethods(iris.MethodOptions)) // <- important for the penlight.
	apiParty.Register(svc)
	apiParty.Handle(new(api.CApi))

	return
}

func getPagination(ctx iris.Context) *model.Pager {
	return &model.Pager{
		Page:     ctx.URLParamInt64Default("page", 1),
		PageSize: ctx.URLParamInt64Default("pagesize", 10),
		UrlPath:  ctx.Path(),
	}
}

func getSiteOptions(svc service.Service, cfg *config.Config) func(ctx iris.Context) (getOption func(name string) string, err error) {
	return func(ctx iris.Context) (func(name string) string, error) {
		options, err := svc.GetSiteOptions()
		if err != nil {
			return nil, err
		}
		ctx.ViewData("options", options)
		path, _ := getDefaultStaticDir(cfg.StaticDir)
		ctx.ViewData("theme",
			fmt.Sprintf("/%s/theme/%s/", path, options["theme"]),
		)
		if err = makeGlobalData(ctx, svc); err != nil {
			return nil, err
		}
		return func(name string) string {
			if value, ok := options[name]; ok {
				return value
			}
			return ""
		}, nil
	}
}

func makeGlobalData(ctx iris.Context, svc service.Service) (err error) {
	ctx.ViewData("hidejs", `<!--[if lt IE 9]>
	<script src="/static/js/html5shiv.min.js"></script>
	<![endif]-->`,
	)
	links, err := svc.GetLinks()
	if err != nil {
		return
	}
	ctx.ViewData("links", links)
	return
}
