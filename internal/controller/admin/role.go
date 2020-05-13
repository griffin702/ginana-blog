package admin

func (c *CAdmin) GetRoleList() (err error) {
	roles, err := c.Svc.GetRoles(c.Pager)
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", roles)
	c.setHeadMetas("角色列表")
	c.Ctx.View("admin/role/list.html")
	return
}
