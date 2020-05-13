package admin

import "ginana-blog/internal/model"

func (c *CAdmin) GetPolicyList() (err error) {
	polices, err := c.Svc.GetPolices(c.Pager)
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", polices)
	c.setHeadMetas("规则列表")
	c.Ctx.View("admin/policy/list.html")
	return
}

func (c *CAdmin) GetPolicyAdd() (err error) {
	c.setHeadMetas("规则创建")
	c.Ctx.View("admin/policy/add.html")
	return
}

func (c *CAdmin) PostPolicyAdd() (err error) {
	req := new(model.CreatePolicyReq)
	if err = c.Ctx.ReadForm(req); err != nil {
		return
	}
	if err = c.Valid(req); err != nil {
		return
	}
	if _, err = c.Svc.CreatePolicy(req); err != nil {
		return
	}
	c.setHeadMetas("规则创建")
	c.ShowMsg("规则已创建")
	return
}

func (c *CAdmin) GetPolicyEditBy(id int64) (err error) {
	policy, err := c.Svc.GetPolicy(id)
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", policy)
	c.setHeadMetas("规则编辑")
	c.Ctx.View("admin/policy/edit.html")
	return
}

func (c *CAdmin) PostPolicyEditBy(id int64) (err error) {
	req := new(model.UpdatePolicyReq)
	if err = c.Ctx.ReadForm(req); err != nil {
		return
	}
	req.ID = id
	if err = c.Valid(req); err != nil {
		return
	}
	if _, err = c.Svc.UpdatePolicy(req); err != nil {
		return
	}
	c.setHeadMetas("规则更新")
	c.ShowMsg("规则已更新")
	return
}

func (c *CAdmin) GetPolicyDeleteBy(id int64) (err error) {
	if err = c.Svc.DeletePolicy(id); err != nil {
		return
	}
	c.setHeadMetas("删除规则")
	c.ShowMsg("规则已删除")
	return
}
