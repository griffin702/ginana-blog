package front

func (c *CFront) Get() {
	c.setHeadMetas()
	c.Ctx.View("front/index.html")
}
