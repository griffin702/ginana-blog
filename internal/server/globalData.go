package server

import (
	"fmt"
	"ginana-blog/internal/config"
	"ginana-blog/internal/model"
	"ginana-blog/internal/service"
	"github.com/kataras/iris/v12"
)

func getPagination(ctx iris.Context) *model.Pager {
	return &model.Pager{
		Page:     ctx.URLParamInt64Default("page", 1),
		PageSize: ctx.URLParamInt64Default("pagesize", 15),
		UrlPath:  ctx.Path(),
	}
}

func getSiteOptions(svc service.Service, cfg *config.Config) model.GetOptionHandler {
	return func(ctx iris.Context) (model.GetOption, error) {
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
	latestArticles, err := svc.GetLatestArticles(5)
	if err != nil {
		return
	}
	ctx.ViewData("latestArticles", latestArticles)
	hotArticles, err := svc.GetHotArticles(5)
	if err != nil {
		return
	}
	ctx.ViewData("hotArticles", hotArticles)
	latestComments, err := svc.GetLatestComments(5)
	if err != nil {
		return
	}
	ctx.ViewData("latestComments", latestComments)
	return
}
