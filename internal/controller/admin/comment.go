package admin

import (
	"ginana-blog/internal/model"
)

func (c *CAdmin) GetCommentList() (err error) {
	comments, err := c.Svc.GetComments(c.Pager, model.CommentQueryParam{Admin: true})
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", comments)
	c.setHeadMetas("用户列表")
	c.Ctx.View("admin/comment/list.html")
	return
}

func (c *CAdmin) PostCommentAdd() {
	req := new(model.CreateCommentReq)
	if err := c.Ctx.ReadJSON(req); err != nil {
		c.Ctx.JSON(c.JsonPlus(false, err))
		return
	}
	req.UserID = c.UserID
	req.IPAddress = c.GetClientIP()
	if err := c.Valid(req); err != nil {
		c.Ctx.JSON(c.JsonPlus(false, err))
		return
	}
	if err := c.Svc.CreateComment(req); err != nil {
		c.Ctx.JSON(c.JsonPlus(false, err))
		return
	}
	c.Ctx.JSON(c.JsonPlus(true, "成功添加评论"))
	return
}

func (c *CAdmin) GetCommentEditBy(id int64) (err error) {
	comment, err := c.Svc.GetComment(id)
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", comment)
	c.setHeadMetas("留言编辑")
	c.Ctx.View("admin/comment/edit.html")
	return
}

func (c *CAdmin) PostCommentEditBy(id int64) (err error) {
	req := new(model.UpdateCommentReq)
	if err = c.Ctx.ReadForm(req); err != nil {
		return
	}
	req.ID = id
	if err = c.Valid(req); err != nil {
		return
	}
	if err = c.Svc.UpdateComment(req); err != nil {
		return
	}
	c.setHeadMetas("留言板")
	c.ShowMsg("留言更新成功")
	return
}

func (c *CAdmin) GetCommentDeleteBy(id int64) (err error) {
	if err = c.Svc.DeleteComment(id); err != nil {
		return
	}
	c.setHeadMetas("评论删除")
	c.ShowMsg("评论已删除")
	return
}
