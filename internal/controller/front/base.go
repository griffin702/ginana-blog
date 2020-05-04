package front

import (
	"ginana-blog/internal/model"
	"ginana-blog/internal/service"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"strings"
)

type CFront struct {
	Ctx          iris.Context
	Session      *sessions.Session
	Svc          service.Service
	GetOption    func(name string) string
	Links        *model.Links
	Pager        *model.Pager
	DisableRight bool
	Hm           service.HelperMap
	Valid        model.Validator
}

func (c *CFront) IsLogin() bool {
	userId := c.Session.Get("userId")
	if userId != nil && userId.(int64) > 0 {
		return true
	}
	return false
}

func (c *CFront) setHeadMetas(params ...string) {
	c.Ctx.ViewData("isLogin", c.IsLogin())
	c.Ctx.ViewData("disableRight", c.DisableRight)
	titleBuf := make([]string, 0, 3)
	if len(params) == 0 && c.GetOption("sitename") != "" {
		titleBuf = append(titleBuf, c.GetOption("sitename"))
	}
	if len(params) > 0 {
		titleBuf = append(titleBuf, params[0])
	}
	titleBuf = append(titleBuf, c.GetOption("subtitle"))
	c.Ctx.ViewData("title", strings.Join(titleBuf, " - "))
	if len(params) > 1 {
		c.Ctx.ViewData("keywords", params[1])
	} else {
		c.Ctx.ViewData("keywords", c.GetOption("keywords"))
	}
	if len(params) > 2 {
		c.Ctx.ViewData("description", params[2])
	} else {
		c.Ctx.ViewData("description", c.GetOption("description"))
	}
}
