package front

import (
	"fmt"
	"ginana-blog/internal/controller"
	"ginana-blog/internal/model"
	"github.com/kataras/iris/v12/mvc"
	"strings"
)

type CFront struct {
	controller.BaseController
}

func (c *CFront) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/about.html", "GetAbout")
	b.Handle("GET", "/life.html", "GetLife")
	b.Handle("GET", "/category.html", "GetCategorys")
	b.Handle("GET", "/mood.html", "GetMoods")
	b.Handle("GET", "/links.html", "GetLinks")
	b.Handle("GET", "/album.html", "GetAlbums")
	b.Handle("GET", "/s/{urlName:path}", "SpecialURL")
}

func (c *CFront) setHeadMetas(params ...string) {
	c.Ctx.ViewData("disableRight", c.DisableRight)
	c.Ctx.ViewData("enableBanner", c.EnableBanner)
	titleBuf := make([]string, 0, 3)
	if len(params) == 0 && c.SiteOptions.SiteName != "" {
		titleBuf = append(titleBuf, c.SiteOptions.SiteName)
	}
	if len(params) > 0 {
		titleBuf = append(titleBuf, params[0])
	}
	titleBuf = append(titleBuf, c.SiteOptions.SubTitle)
	c.Ctx.ViewData("title", strings.Join(titleBuf, " - "))
	if len(params) > 1 {
		c.Ctx.ViewData("keywords", params[1])
	} else {
		c.Ctx.ViewData("keywords", c.SiteOptions.Keywords)
	}
	if len(params) > 2 {
		c.Ctx.ViewData("description", params[2])
	} else {
		c.Ctx.ViewData("description", c.SiteOptions.Description)
	}
}

func (c *CFront) Get() (err error) {
	tags, err := c.Svc.GetTagsLimit6()
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", tags)
	c.EnableBanner = true
	c.setHeadMetas("首页")
	c.Ctx.View("front/index.html")
	return
}

func (c *CFront) GetAbout() (err error) {
	c.DisableRight = true
	c.setHeadMetas("关于我")
	c.Ctx.View("front/about.html")
	return
}

func (c *CFront) GetArticleBy(id int64) (err error) {
	article, err := c.Svc.GetArticle(id)
	if err != nil {
		return
	}
	if err = c.Svc.AddViews(article); err != nil {
		return
	}
	c.Ctx.ViewData("data", article)
	comments, err := c.Svc.GetComments(c.Pager, model.CommentQueryParam{ArticleID: id})
	if err != nil {
		return
	}
	c.Ctx.ViewData("comments", comments)
	c.setHeadMetas(article.Title)
	c.Ctx.View("front/article.html")
	return
}

func (c *CFront) GetLife() (err error) {
	articles, err := c.Svc.GetArticles(c.Pager)
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", articles)
	c.setHeadMetas("成长录")
	c.Ctx.View("front/life.html")
	return
}

func (c *CFront) GetCategorys() (err error) {
	c.Pager.PageSize = 50
	tags, err := c.Svc.GetTags(c.Pager)
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", tags)
	c.setHeadMetas("归类归档")
	c.Ctx.View("front/category.html")
	return
}

func (c *CFront) GetCategoryBy(id int64) (err error) {
	articles, err := c.Svc.GetArticles(c.Pager, model.ArticleQueryParam{TagID: id})
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", articles)
	c.setHeadMetas("归类归档")
	c.Ctx.View("front/categoryList.html")
	return
}

func (c *CFront) GetMoods() (err error) {
	moods, err := c.Svc.GetMoods(c.Pager)
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", moods)
	c.DisableRight = true
	c.setHeadMetas("碎言碎语")
	c.Ctx.View("front/mood.html")
	return
}

func (c *CFront) GetLinks() (err error) {
	comments, err := c.Svc.GetComments(c.Pager, model.CommentQueryParam{
		ArticleID: 0,
	})
	if err != nil {
		return
	}
	c.Ctx.ViewData("comments", comments)
	c.setHeadMetas("友情链接")
	c.Ctx.View("front/links.html")
	return
}

func (c *CFront) GetAlbums() (err error) {
	albums, err := c.Svc.GetAlbums(c.Pager)
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", albums)
	c.DisableRight = true
	c.setHeadMetas("光影瞬间")
	c.Ctx.View("front/album.html")
	return
}

func (c *CFront) GetAlbumBy(id int64) (err error) {
	album, err := c.Svc.GetAlbum(id)
	if err != nil {
		return
	}
	photos, err := c.Svc.GetPhotos(c.Pager, id)
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", photos)
	c.DisableRight = true
	c.setHeadMetas(fmt.Sprintf("相册 %s 内的照片", album.Name))
	c.Ctx.View("front/photo.html")
	return
}

func (c *CFront) SpecialURL() (err error) {
	urlName := c.Ctx.Params().GetStringDefault("urlName", "")
	if urlName == "" {
		return c.Hm.GetMessage(404, "404 not found")
	}
	article, err := c.Svc.GetArticleByUrlName(urlName)
	if err != nil {
		return c.Hm.GetMessage(404, "404 not found")
	}
	c.Ctx.ViewData("data", article)
	comments, err := c.Svc.GetComments(c.Pager, model.CommentQueryParam{ArticleID: article.ID})
	if err != nil {
		return
	}
	c.Ctx.ViewData("comments", comments)
	c.setHeadMetas(article.Title)
	c.Ctx.View("front/article.html")
	return
}
