package admin

import "ginana-blog/internal/model"

func (c *CAdmin) GetMoodList() (err error) {
	moods, err := c.Svc.GetMoods(c.Pager)
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", moods)
	c.setHeadMetas("心情列表")
	c.Ctx.View("admin/mood/list.html")
	return
}

func (c *CAdmin) GetMoodAdd() (err error) {
	c.setHeadMetas("添加心情")
	c.Ctx.View("admin/mood/add.html")
	return
}

func (c *CAdmin) PostMoodAdd() (err error) {
	req := new(model.MoodReq)
	if err = c.Ctx.ReadForm(req); err != nil {
		return
	}
	if err = c.Valid(req); err != nil {
		return
	}
	if err = c.Svc.CreateMood(req); err != nil {
		return
	}
	c.setHeadMetas("添加心情")
	c.ShowMsg("添加心情成功")
	return
}
