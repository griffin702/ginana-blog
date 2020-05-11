package admin

import "ginana-blog/internal/model"

func (c *CAdmin) GetTagList() (err error) {
	tags, err := c.Svc.GetTags(c.Pager, model.TagQueryParam{Admin: true})
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", tags)
	c.setHeadMetas("标签列表")
	c.Ctx.View("admin/tag/list.html")
	return
}

func (c *CAdmin) PostTagList() (err error) {
	req := new(model.TagListReq)
	if err = c.Ctx.ReadForm(req); err != nil {
		return
	}
	if err = c.Valid(req); err != nil {
		return
	}
	if err = c.Svc.BatchTag(req); err != nil {
		return
	}
	c.setHeadMetas("标签批量处理")
	c.ShowMsg("标签批量处理完成")
	return
}
