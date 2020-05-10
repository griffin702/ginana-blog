package admin

import (
	"ginana-blog/internal/model"
)

func (c *CAdmin) GetUserList() (err error) {
	users, err := c.Svc.GetUsers(c.Pager)
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", users)
	c.setHeadMetas("用户列表")
	c.Ctx.View("admin/user/list.html")
	return
}

func (c *CAdmin) GetUserAdd() (err error) {
	c.setHeadMetas("用户创建")
	c.Ctx.View("admin/user/add.html")
	return
}

func (c *CAdmin) PostUserAdd() (err error) {
	req := new(model.CreateUserReq)
	if err = c.Ctx.ReadForm(req); err != nil {
		return
	}
	if req.Nickname == "" {
		req.Nickname = req.Username
	}
	if err = c.Valid(req); err != nil {
		return
	}
	if _, err = c.Svc.CreateUser(req); err != nil {
		return
	}
	c.setHeadMetas("用户创建")
	c.ShowMsg("用户已创建")
	return
}

func (c *CAdmin) GetUserEditBy(id int64) (err error) {
	user, err := c.Svc.GetUser(id)
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", user)
	c.setHeadMetas("用户编辑")
	c.Ctx.View("admin/user/edit.html")
	return
}

func (c *CAdmin) PostUserEditBy(id int64) (err error) {
	req := new(model.UpdateUserReq)
	req.ID = id
	if err = c.Ctx.ReadForm(req); err != nil {
		return
	}
	if err = c.Valid(req); err != nil {
		return
	}
	if _, err = c.Svc.UpdateUser(req); err != nil {
		return
	}
	c.setHeadMetas("用户更新")
	c.ShowMsg("用户已更新")
	return
}

func (c *CAdmin) GetUserDeleteBy(id int64) (err error) {
	if err = c.Svc.DeleteUser(id); err != nil {
		return
	}
	c.setHeadMetas("用户删除")
	c.ShowMsg("用户已删除")
	return
}
