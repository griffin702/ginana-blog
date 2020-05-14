package admin

import (
	"ginana-blog/internal/model"
)

func (c *CAdmin) GetLinkList() (err error) {
	c.setHeadMetas("友链列表")
	c.Ctx.View("admin/link/list.html")
	return
}

func (c *CAdmin) GetLinkAdd() (err error) {
	c.setHeadMetas("友链创建")
	c.Ctx.View("admin/link/add.html")
	return
}

func (c *CAdmin) PostLinkAdd() (err error) {
	req := new(model.CreateLinkReq)
	if err = c.Ctx.ReadForm(req); err != nil {
		return
	}
	if err = c.Valid(req); err != nil {
		return
	}
	if _, err = c.Svc.CreateLink(req); err != nil {
		return
	}
	c.setHeadMetas("友链创建")
	c.ShowMsg("友链已创建")
	return
}

func (c *CAdmin) GetLinkEditBy(id int64) (err error) {
	link, err := c.Svc.GetLink(id)
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", link)
	c.setHeadMetas("友链编辑")
	c.Ctx.View("admin/link/edit.html")
	return
}

func (c *CAdmin) PostLinkEditBy(id int64) (err error) {
	req := new(model.UpdateLinkReq)
	if err = c.Ctx.ReadForm(req); err != nil {
		return
	}
	req.ID = id
	if err = c.Valid(req); err != nil {
		return
	}
	if _, err = c.Svc.UpdateLink(req); err != nil {
		return
	}
	c.setHeadMetas("友链更新")
	c.ShowMsg("友链已更新")
	return
}

func (c *CAdmin) GetLinkDeleteBy(id int64) (err error) {
	if err = c.Svc.DeleteLink(id); err != nil {
		return
	}
	c.setHeadMetas("删除友链")
	c.ShowMsg("友链已删除")
	return
}
