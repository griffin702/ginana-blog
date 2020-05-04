package admin

func (c *CAdmin) Get() {
	c.Ctx.View("admin/index.html")
}
