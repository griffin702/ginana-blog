package front

import (
	"github.com/kataras/iris/v12/mvc"
)

func (c *CFront) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/about.html", "GetAbout")
	b.Handle("GET", "/life.html", "GetLife")
	b.Handle("GET", "/category.html", "GetTags")
	b.Handle("GET", "/mood.html", "GetMoods")
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

func (c *CFront) GetTags() (err error) {
	resp, err := c.Svc.GetTags()
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", resp)
	c.setHeadMetas("归类归档")
	c.Ctx.View("front/category.html")
	return
}

func (c *CFront) GetMoods() (err error) {
	resp, err := c.Svc.GetMoods(c.Pager)
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", resp)
	c.DisableRight = true
	c.setHeadMetas("碎言碎语")
	c.Ctx.View("front/mood.html")
	return
}
