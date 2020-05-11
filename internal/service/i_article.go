package service

import (
	"ginana-blog/internal/model"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"strings"
)

func (s *service) GetArticles(p *model.Pager, prs ...model.ArticleQueryParam) (res *model.Articles, err error) {
	var pr model.ArticleQueryParam
	if len(prs) > 0 {
		pr = prs[0]
	}
	if pr.Order == "" {
		pr.Order = "istop desc, id desc"
	}
	res = new(model.Articles)
	query := s.db.Model(&res.List)
	query.Where("status = 0").Count(&res.CountStatus0)
	query.Where("status = 1").Count(&res.CountStatus1)
	query.Where("status = 2").Count(&res.CountStatus2)
	if pr.Search != "" && pr.Keyword != "" { // 搜索
		switch pr.Search {
		case "title":
			query = query.Where("title like ?", "%"+pr.Keyword+"%")
		case "author":
			var userIdList []int64
			s.db.Model(&model.User{}).Select("id").
				Where("nickname like ?", "%"+pr.Keyword+"%").Pluck("id", &userIdList)
			query = query.Having("user_id in (?)", userIdList)
		case "tag":
			var tagList []*model.Tag
			s.db.Model(&tagList).Preload("Articles").Find(&tagList, "name like ?", "%"+pr.Keyword+"%")
			var idList []int64
			for _, tag := range tagList {
				for _, art := range tag.Articles {
					idList = append(idList, art.ID)
				}
			}
			query = query.Having("id in (?)", idList)
		}
	}
	if pr.TagID > 0 {
		query = query.Joins("left join w_article_tags ON w_article_tags.article_id = w_article.id "+
			"and w_article_tags.tag_id = ?", pr.TagID)
	}
	query = query.Where("status = ?", pr.Status)
	query.Count(&p.AllCount)
	query = query.Order(pr.Order).Preload("User").Preload("Tags")
	if err = query.Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize).Find(&res.List).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	res.Status = pr.Status
	res.Search = pr.Search
	res.Keyword = pr.Keyword
	res.Pager = p
	return
}

func (s *service) GetArticle(id int64) (article *model.Article, err error) {
	article = new(model.Article)
	article.ID = id
	if err = s.db.Model(article).Preload("User").Preload("Tags").
		Find(article).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	var prev, next model.Article
	err = s.db.Model(&prev).Last(&prev, "id < ?", id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, s.hm.GetMessage(1001, err)
	}
	if err == nil {
		article.Prev = &prev
	}
	err = s.db.Model(&next).First(&next, "id > ?", id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, s.hm.GetMessage(1001, err)
	}
	if err == nil {
		article.Next = &next
	}
	return article, nil
}

func (s *service) GetArticleByUrlName(urlName string) (article *model.Article, err error) {
	article = new(model.Article)
	if err = s.db.Model(article).Preload("User").Preload("Tags").
		Find(article, "urlname = ?", urlName).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	var prev, next model.Article
	err = s.db.Model(&prev).Last(&prev, "id < ?", article.ID).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, s.hm.GetMessage(1001, err)
	}
	if err == nil {
		article.Prev = &prev
	}
	err = s.db.Model(&next).First(&next, "id > ?", article.ID).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, s.hm.GetMessage(1001, err)
	}
	if err == nil {
		article.Next = &next
	}
	return article, nil
}

func (s *service) GetLatestArticles(limit int) (articles []*model.Article, err error) {
	key := s.hm.GetCacheKey(4)
	err = s.mc.Get(key, &articles)
	if err != nil {
		if err = s.db.Model(&articles).Order("created_at desc").Limit(limit).Find(&articles).Error; err != nil {
			return nil, s.hm.GetMessage(1001, err)
		}
		if err = s.mc.Set(key, &articles); err != nil {
			return nil, s.hm.GetMessage(1002, err)
		}
	}
	return
}

func (s *service) GetHotArticles(limit int) (articles []*model.Article, err error) {
	key := s.hm.GetCacheKey(5)
	err = s.mc.Get(key, &articles)
	if err != nil {
		if err = s.db.Model(&articles).Order("views desc").Limit(limit).Find(&articles).Error; err != nil {
			return nil, s.hm.GetMessage(1001, err)
		}
		if err = s.mc.Set(key, &articles); err != nil {
			return nil, s.hm.GetMessage(1002, err)
		}
	}
	return
}

func (s *service) CountArticles() (count int64) {
	article := new(model.Article)
	s.db.Model(article).Count(&count)
	return
}

func (s *service) CreateArticle(req *model.ArticleReq) (article *model.Article, err error) {
	article = new(model.Article)
	article.Title = req.Title
	article.Color = req.Color
	article.Urlname = req.Urlname
	article.Urltype = req.Urltype
	article.Istop = req.Istop
	article.Status = req.Status
	article.Content = req.ContentMarkdownDoc
	article.Cover = req.Cover
	article.UserID = req.UserID
	tags := strings.Split(req.Tags, ",")
	for _, name := range tags {
		tag, err := s.GetTagByName(name)
		if err == gorm.ErrRecordNotFound {
			tag = new(model.Tag)
			tag.Name = name
			err = nil
		} else if err != nil {
			return nil, s.hm.GetMessage(500, err)
		}
		article.Tags = append(article.Tags, tag)
	}
	if err = s.db.Create(article).Error; err != nil {
		return nil, s.hm.GetMessage(1002, err)
	}
	s.deleteArticleCache()
	return
}

func (s *service) UpdateArticle(req *model.ArticleReq) (article *model.Article, err error) {
	article = new(model.Article)
	if err = s.db.Find(article, "id = ?", req.ID).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	article.Title = req.Title
	article.Color = req.Color
	article.Urlname = req.Urlname
	article.Urltype = req.Urltype
	article.Istop = req.Istop
	article.Status = req.Status
	article.Content = req.ContentMarkdownDoc
	article.Cover = req.Cover
	m, err := s.tool.StructToMap(article)
	if err != nil {
		return nil, s.hm.GetMessage(500, err)
	}
	tx := s.db.Begin()
	if err = tx.Model(article).Update(m).Error; err != nil {
		tx.Rollback()
		return nil, s.hm.GetMessage(1003, err)
	}
	tags := strings.Split(req.Tags, ",")
	for _, name := range tags {
		tag, err := s.GetTagByName(name)
		if err == gorm.ErrRecordNotFound {
			tag = new(model.Tag)
			tag.Name = name
			err = nil
		} else if err != nil {
			return nil, s.hm.GetMessage(500, err)
		}
		article.Tags = append(article.Tags, tag)
	}
	if err = tx.Model(article).Update(article).Error; err != nil {
		tx.Rollback()
		return nil, s.hm.GetMessage(1003, err)
	}
	tx.Commit()
	s.deleteArticleCache()
	return
}

func (s *service) DeleteArticle(id int64) (err error) {
	article := new(model.Article)
	if err = s.db.Delete(article, "id = ?", id).Error; err != nil {
		return s.hm.GetMessage(1004, err)
	}
	return
}

func (s *service) BatchArticle(req *model.ArticleListReq) (err error) {
	article := new(model.Article)
	switch req.Option {
	case "public":
		err = s.db.Model(article).Where("id in (?)", req.IDs).Update("status", 0).Error
	case "template":
		err = s.db.Model(article).Where("id in (?)", req.IDs).Update("status", 1).Error
	case "recycle":
		err = s.db.Model(article).Where("id in (?)", req.IDs).Update("status", 2).Error
	case "delete":
		if err = s.db.Model(article).Where("id in (?)", req.IDs).Delete(article).Error; err != nil {
			return s.hm.GetMessage(1004, err)
		}
	}
	if err != nil {
		return s.hm.GetMessage(1003, err)
	}
	s.deleteArticleCache()
	return
}

func (s *service) PushBaiDu(url string) (string, error) {
	resp, err := http.Post("http://data.zz.baidu.com/urls?site=https://www.inana.top&token=d0Dca7O4TosN7655",
		"application/x-www-form-urlencoded",
		strings.NewReader(url))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (s *service) deleteArticleCache() {
	s.mc.Delete(s.hm.GetCacheKey(4))
	s.mc.Delete(s.hm.GetCacheKey(5))
}
