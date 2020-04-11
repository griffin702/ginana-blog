package router

import (
	"ginana-blog/internal/config"
	"ginana-blog/internal/server/resp"
	"ginana-blog/library/ecode"
	"ginana-blog/library/log"
	"ginana-blog/library/tools"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"html/template"
	"strings"
	"time"
)

func NewIris(cfg *config.Config) (e *iris.Application) {
	e = iris.New()
	golog.Install(log.GetLogger())
	customLogger := logger.New(logger.Config{
		Status: true, IP: true, Method: true, Path: true, Query: true,
	})
	e.Use(customLogger, recover.New())
	e.Logger().SetLevel(cfg.IrisLogLevel)
	e.Use(func(ctx iris.Context) {
		ctx.Gzip(cfg.EnableGzip)
		ctx.Next()
	})
	initTemplate(e, cfg)
	initStaticDir(e, cfg)
	e.OnAnyErrorCode(customLogger, func(ctx iris.Context) {
		ctx.JSON(resp.JsonError(ecode.Errorf(ctx.GetStatusCode())))
	})
	//// Cors
	//e.Use(mdw.CORS([]string{"*"}))
	//// Swagger
	//handle := mdw.SwaggerHandler("http://127.0.0.1:8000/swagger/doc.json")
	//e.GET("/swagger/*any", handle)
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
}

func initStaticDir(e *iris.Application, cfg *config.Config) {
	if !cfg.EnableTemplate {
		return
	}
	staticDirList := strings.Split(cfg.StaticDir, " ")
	for _, v := range staticDirList {
		path := strings.Split(v, ":")
		if len(path) == 2 {
			e.HandleDir(path[0], path[1], iris.DirOptions{Gzip: true})
		}
	}
}

// template function
func dateFormat(t time.Time, format string) (template.HTML, error) {
	return template.HTML(tools.New().TimeFormat(&t, format)), nil
}

func str2html(str string) (template.HTML, error) {
	return template.HTML(str), nil
}
