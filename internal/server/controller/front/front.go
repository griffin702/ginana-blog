package front

func (c *CFront) Get() (err error) {
	c.setHeadMetas("首页")
	c.Ctx.View("front/index.html")
	return
}
