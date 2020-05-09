package admin

import (
	"ginana-blog/internal/model"
)

func (c *CAdmin) GetArticleList() (err error) {
	status := c.Ctx.URLParamIntDefault("status", 0)
	articles, err := c.Svc.GetArticles(c.Pager, model.ArticleQueryParam{Status: status})
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", articles)
	c.setHeadMetas("文章列表")
	c.Ctx.View("admin/article/list.html")
	return
}

func (c *CAdmin) GetArticleAdd() (err error) {
	c.setHeadMetas("文章创建")
	c.Ctx.View("admin/article/add.html")
	return
}
