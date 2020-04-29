package front

import (
	"github.com/kataras/iris/v12/mvc"
)

func (c *CFront) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/about.html", "GetAbout")
	b.Handle("GET", "/life.html", "GetLife")
}

func (c *CFront) Get() (err error) {
	resp, err := c.Svc.GetArticles(c.Pager)
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", resp)
	c.setHeadMetas("首页")
	c.Ctx.View("front/index.html")
	return
}

func (c *CFront) GetAbout() (err error) {
	c.DisableRight = true
	c.setHeadMetas("关于我")
	c.Ctx.View("front/about.html")
	return
}

func (c *CFront) GetLife() (err error) {
	resp, err := c.Svc.GetArticles(c.Pager)
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", resp)
	c.setHeadMetas("成长录")
	c.Ctx.View("front/life.html")
	return
}
