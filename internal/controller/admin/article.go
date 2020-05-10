package admin

import (
	"ginana-blog/internal/model"
)

func (c *CAdmin) GetArticleList() (err error) {
	status := c.Ctx.URLParamIntDefault("status", 0)
	search := c.Ctx.URLParam("search")
	keyword := c.Ctx.URLParam("keyword")
	articles, err := c.Svc.GetArticles(c.Pager, model.ArticleQueryParam{
		Status:  status,
		Search:  search,
		Keyword: keyword,
	})
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", articles)
	c.setHeadMetas("文章列表")
	c.Ctx.View("admin/article/list.html")
	return
}

func (c *CAdmin) PostArticleList() (err error) {
	req := new(model.ArticleListReq)
	if err = c.Ctx.ReadForm(req); err != nil {
		return
	}
	if err = c.Valid(req); err != nil {
		return
	}
	if err = c.Svc.BatchArticle(req); err != nil {
		return
	}
	c.setHeadMetas("文章批量处理")
	c.ShowMsg("文章批量处理完成")
	return
}

func (c *CAdmin) GetArticleAdd() (err error) {
	c.setHeadMetas("文章创建")
	c.Ctx.View("admin/article/add.html")
	return
}

func (c *CAdmin) PostArticleAdd() (err error) {
	req := new(model.ArticleReq)
	if err = c.Ctx.ReadForm(req); err != nil {
		return
	}
	req.UserID = c.UserID
	if err = c.Valid(req); err != nil {
		return
	}
	if _, err = c.Svc.CreateArticle(req); err != nil {
		return
	}
	c.setHeadMetas("文章创建")
	c.ShowMsg("文章已创建")
	return
}

func (c *CAdmin) GetArticleEditBy(id int64) (err error) {
	article, err := c.Svc.GetArticle(id)
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", article)
	c.setHeadMetas("文章创建")
	c.Ctx.View("admin/article/edit.html")
	return
}

func (c *CAdmin) PostArticleEditBy(id int64) (err error) {
	req := new(model.ArticleReq)
	if err = c.Ctx.ReadForm(req); err != nil {
		return
	}
	req.ID = id
	if err = c.Valid(req); err != nil {
		return
	}
	if _, err = c.Svc.UpdateArticle(req); err != nil {
		return
	}
	c.setHeadMetas("文章更新")
	c.ShowMsg("文章已更新")
	return
}

func (c *CAdmin) GetArticleDeleteBy(id int64) (err error) {
	if err = c.Svc.DeleteArticle(id); err != nil {
		return
	}
	c.setHeadMetas("文章删除")
	c.ShowMsg("文章已删除")
	return
}
