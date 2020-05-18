package server

import (
	"fmt"
	"ginana-blog/internal/config"
	"ginana-blog/internal/model"
	"ginana-blog/internal/service"
	"github.com/griffin702/ginana/library/ecode"
	"github.com/griffin702/ginana/library/log"
	"github.com/griffin702/service/tools"
	"github.com/kataras/iris/v12"
	"strconv"
	"strings"
)

func getPagination(svc service.Service) model.GetPagination {
	return func(ctx iris.Context) (p *model.Pager, err error) {
		options, err := svc.GetSiteOptions()
		if err != nil {
			return nil, err
		}
		p = new(model.Pager)
		p.Page = ctx.URLParamInt64Default("page", 1)
		size, err := strconv.Atoi(options.PageSize)
		if err != nil {
			return nil, err
		}
		p.PageSize = ctx.URLParamInt64Default("pagesize", int64(size))
		p.UrlPath = ctx.Path()
		p.UrlParams = ctx.URLParams()
		return
	}
}

func getSiteOptions(svc service.Service, cfg *config.Config) model.OptionHandler {
	return func(ctx iris.Context) (*model.Option, error) {
		options, err := svc.GetSiteOptions()
		if err != nil {
			return nil, err
		}
		ctx.ViewData("options", options)
		path, _ := getDefaultStaticDir(cfg.StaticDir)
		ctx.ViewData("theme",
			fmt.Sprintf("/%s/theme/%s/", path, options.Theme),
		)
		if err = makeGlobalData(ctx, svc); err != nil {
			return nil, err
		}
		return options, nil
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

func getClientIP(ctx iris.Context) model.GetClientIP {
	return func() string {
		s := ctx.GetHeader("X-Real-IP")
		if s == "" {
			forwarded := ctx.GetHeader("X-Forwarded-For")
			if forwarded != "" {
				list := strings.Split(forwarded, ":")
				if len(list) > 0 {
					s = list[0]
				}
			} else {
				s = strings.Split(ctx.RemoteAddr(), ":")[0]
			}
		}
		return s
	}
}

func getTools(_ iris.Context) *tools.Tool {
	return tools.Tools
}

func getConfigs(_ iris.Context) *config.Config {
	return config.Global()
}

func jsonPlus(_ iris.Context) model.JsonPlus {
	return func(data interface{}, msg interface{}) *model.JSON {
		ec := ecode.Cause(msg)
		return &model.JSON{
			Code:    ec.Code(),
			Message: ec.Message(),
			Data:    data,
		}
	}
}

func errorHandler(ctx iris.Context, err error) {
	jp := jsonPlus(ctx)(nil, err)
	log.Errorf("%d, %s", jp.Code, jp.Message)
	redirect := ctx.GetReferrer().Path
	if redirect == "" {
		redirect = "/"
	}
	ctx.ViewData("redirect", redirect)
	ctx.ViewData("error", jp)
	ctx.View("message/error.html")
}
