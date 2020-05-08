package admin

func (c *CAdmin) GetAccountInfo() (err error) {
	user, err := c.Svc.GetUser(c.UserID)
	if err != nil {
		return
	}
	c.Ctx.ViewData("user", user)
	c.setHeadMetas("个人资料")
	c.Ctx.View("admin/account/info.html")
	return
}
