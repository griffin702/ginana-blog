package admin

import (
	"ginana-blog/internal/model"
)

func (c *CAdmin) PostComments() (err error) {
	req := new(model.Comment)
	if err = c.Ctx.ReadForm(req); err != nil {
		return
	}
	req.UserID = c.UserID
	req.IPAddress = c.GetClientIP()
	if err = c.Valid(req); err != nil {
		return
	}
	if err = c.Svc.PostComment(req); err != nil {
		return
	}
	c.setHeadMetas("留言板")
	c.ShowMsg("留言成功")
	return
}
