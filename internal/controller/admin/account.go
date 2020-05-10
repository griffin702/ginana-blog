package admin

import "ginana-blog/internal/model"

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

func (c *CAdmin) PostAccountInfo() (err error) {
	req := new(model.UpdateUserReq)
	if err = c.Ctx.ReadForm(req); err != nil {
		return
	}
	req.ID = c.UserID
	if err = c.Valid(req); err != nil {
		return
	}
	if _, err = c.Svc.UpdateAccount(req); err != nil {
		return
	}
	c.setHeadMetas("修改当前用户资料")
	c.ShowMsg("当前用户资料已更新")
	return
}
