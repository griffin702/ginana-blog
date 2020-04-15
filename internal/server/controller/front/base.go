package front

import (
	"ginana-blog/internal/model"
	"ginana-blog/internal/service"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"strings"
)

type CFront struct {
	Ctx     iris.Context
	Session *sessions.Session
	Svc     service.Service
	Page    *model.Pagination
}

func (c *CFront) IsLogin() bool {
	userId := c.Session.Get("userId")
	if userId != nil && userId.(int64) > 0 {
		return true
	}
	return false
}

func (c *CFront) setHeadMetas(params ...string) {
	c.Ctx.ViewData("IsLogin", c.IsLogin())
	titleBuf := make([]string, 0, 3)
	options, _ := c.Svc.GetSiteOptions()
	if len(params) == 0 && options["sitename"] != "" {
		titleBuf = append(titleBuf, options["sitename"])
	}
	if len(params) > 0 {
		titleBuf = append(titleBuf, params[0])
	}
	titleBuf = append(titleBuf, options["subtitle"])
	c.Ctx.ViewData("title", strings.Join(titleBuf, " - "))
	if len(params) > 1 {
		c.Ctx.ViewData("keywords", params[1])
	} else {
		c.Ctx.ViewData("keywords", options["keywords"])
	}
	if len(params) > 2 {
		c.Ctx.ViewData("description", params[2])
	} else {
		c.Ctx.ViewData("description", options["description"])
	}
}
