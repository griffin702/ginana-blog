package server

import (
	"ginana-blog/internal/config"
	"ginana-blog/internal/service"
	"github.com/griffin702/service/tools"
	"github.com/kataras/iris/v12"
	"html/template"
	"strings"
	"time"
)

func initTemplate(e *iris.Application, svc service.Service, cfg *config.Config) {
	if !cfg.EnableTemplate {
		return
	}
	tmpl := iris.HTML(cfg.ViewsPath, ".html").
		Reload(cfg.ReloadTemplate)
	tmpl.AddFunc("date", dateFormat)
	tmpl.AddFunc("str2html", str2html)
	tmpl.AddFunc("permission", permission(svc))
	e.RegisterView(tmpl)
	return
}

// template function
func dateFormat(t time.Time, format string) (template.HTML, error) {
	return template.HTML(tools.New().TimeFormat(t, format)), nil
}

func str2html(str string) (template.HTML, error) {
	return template.HTML(str), nil
}

func permission(svc service.Service) func(int64, string, string) bool {
	return func(userId int64, router, method string) (isAuth bool) {
		return svc.CheckPermission(userId, router, strings.ToUpper(method))
	}
}
