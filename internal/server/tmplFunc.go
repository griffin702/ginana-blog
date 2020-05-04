package server

import (
	"ginana-blog/internal/config"
	"github.com/griffin702/ginana/library/tools"
	"github.com/kataras/iris/v12"
	"html/template"
	"time"
)

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

// template function
func dateFormat(t time.Time, format string) (template.HTML, error) {
	return template.HTML(tools.New().TimeFormat(&t, format)), nil
}

func str2html(str string) (template.HTML, error) {
	return template.HTML(str), nil
}
