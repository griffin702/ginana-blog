package front

func (c *CFront) Get() {
	c.setHeadMetas("首页")
	c.Ctx.View("front/index.html")
}
