package admin

import "ginana-blog/internal/model"

func (c *CAdmin) GetGocmdPhone_listList() (err error) {
	phoneList, err := c.Svc.GetPhoneLists(c.Pager)
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", phoneList)
	c.setHeadMetas("手机列表")
	c.Ctx.View("admin/gocmd/phone_list/list.html")
	return
}

func (c *CAdmin) GetGocmdPhone_listAdd() (err error) {
	c.setHeadMetas("添加手机列表")
	c.Ctx.View("admin/gocmd/phone_list/add.html")
	return
}

func (c *CAdmin) PostGocmdPhone_listAdd() (err error) {
	req := new(model.CreatePhoneListReq)
	if err = c.Ctx.ReadForm(req); err != nil {
		return
	}
	if err = c.Valid(req); err != nil {
		return
	}
	if _, err = c.Svc.CreatePhoneList(req); err != nil {
		return
	}
	c.setHeadMetas("添加手机列表")
	c.ShowMsg("添加手机列表成功")
	return
}

func (c *CAdmin) GetGocmdPhone_listEditBy(id int64) (err error) {
	phoneList, err := c.Svc.GetPhoneList(id)
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", phoneList)
	c.setHeadMetas("手机列表编辑")
	c.Ctx.View("admin/gocmd/phone_list/edit.html")
	return
}

func (c *CAdmin) PostGocmdPhone_listEditBy(id int64) (err error) {
	req := new(model.UpdatePhoneListReq)
	if err = c.Ctx.ReadForm(req); err != nil {
		return
	}
	req.ID = id
	if err = c.Valid(req); err != nil {
		return
	}
	if _, err = c.Svc.UpdatePhoneList(req); err != nil {
		return
	}
	c.setHeadMetas("手机列表更新")
	c.ShowMsg("手机列表已更新")
	return
}
