package router

import (
	"fmt"
	"ginana-blog/internal/config"
	"ginana-blog/internal/server/resp"
	"ginana-blog/internal/service"
	"ginana-blog/library/ecode"
	"ginana-blog/library/log"
	"ginana-blog/library/mdw"
	"ginana-blog/library/tools"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"html/template"
	"strings"
	"time"
)

func NewIris(svc service.Service, cfg *config.Config) (e *iris.Application) {
	e = iris.New()
	//e.Use(iris.Cache304(10 * time.Second))
	golog.Install(log.GetLogger())
	customLogger := logger.New(logger.Config{
		Status: true, IP: true, Method: true, Path: true, Query: true,
		//MessageHeaderKeys: []string{"User-Agent"},
	})
	e.Use(customLogger, recover.New())
	e.Logger().SetLevel(cfg.IrisLogLevel)
	initTemplate(e, cfg)
	initStaticDir(e, cfg)
	e.Use(func(ctx iris.Context) {
		ctx.Gzip(cfg.EnableGzip)
		ctx.Next()
	})
	e.OnAnyErrorCode(customLogger, func(ctx iris.Context) {
		ctx.JSON(resp.PlusJson(nil, ecode.Errorf(ctx.GetStatusCode())))
	})
	e.UseGlobal(globalData(svc, cfg))
	// Swagger
	handle := mdw.SwaggerHandler("http://127.0.0.1:8000/swagger/doc.json")
	e.Get("/swagger/*any", handle)
	return
}

func initTemplate(e *iris.Application, cfg *config.Config) {
	if !cfg.EnableTemplate {
		return
	}
	tmpl := iris.HTML(cfg.ViewsPath, ".html").
		Reload(cfg.ReloadTemplate)
	tmpl.AddFunc("date", dateFormat)
	tmpl.AddFunc("str2html", str2html)
	e.RegisterView(tmpl)
	return
}

func initStaticDir(e *iris.Application, cfg *config.Config) {
	if !cfg.EnableTemplate {
		return
	}
	staticDirList := strings.Split(cfg.StaticDir, " ")
	if len(staticDirList) > 0 {
		path := strings.Split(staticDirList[0], ":")
		e.Favicon(fmt.Sprintf("%s/favicon.ico", path[1]))
	}
	for _, v := range staticDirList {
		path := strings.Split(v, ":")
		if len(path) == 2 {
			e.HandleDir(path[0], path[1], iris.DirOptions{Gzip: true})
		}
	}
	return
}

func getDefaultStaticDir(conf string) (path, dir string) {
	staticDirList := strings.Split(conf, " ")
	if len(staticDirList) > 0 {
		def := strings.Split(staticDirList[0], ":")
		if len(def) == 2 {
			return def[0], def[1]
		}
	}
	return
}

// template function
func dateFormat(t time.Time, format string) (template.HTML, error) {
	return template.HTML(tools.New().TimeFormat(&t, format)), nil
}

func str2html(str string) (template.HTML, error) {
	return template.HTML(str), nil
}

func globalData(svc service.Service, cfg *config.Config) iris.Handler {
	return func(ctx iris.Context) {
		res, err := svc.GetSiteOptions()
		if err != nil {
			ctx.JSON(resp.PlusJson(nil, err))
			ctx.StopExecution()
			return
		}
		//ctx.ViewData("options", res)
		path, _ := getDefaultStaticDir(cfg.StaticDir)
		ctx.ViewData("theme",
			fmt.Sprintf("/%s/theme/%s/", path, res["theme"]),
		)
		ctx.ViewData("hidejs", `<!--[if lt IE 9]>
	<script src="/static/js/html5shiv.min.js"></script>
	<![endif]-->`,
		)
		ctx.Next()
	}
}
