package front

//func (c *CFront) BeforeActivation(b mvc.BeforeActivation) {
//
//}

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
