package service

import (
	"ginana-blog/internal/model"
	"github.com/jinzhu/gorm"
)

func (s *service) GetArticles(p *model.Pager, prs ...model.ArticleQueryParam) (res *model.Articles, err error) {
	var pr model.ArticleQueryParam
	if len(prs) > 0 {
		pr = prs[0]
	}
	if pr.Order == "" {
		pr.Order = "id"
	}
	res = new(model.Articles)
	query := s.db.Model(&res.List)
	if pr.TagID > 0 {
		query = query.Joins("left join w_article_tags ON w_article_tags.article_id = w_article.id where w_article_tags.tag_id = ?", pr.TagID)
	}
	query.Count(&p.AllCount)
	query = query.Order(pr.Order).Preload("User").Preload("Tags")
	if err = query.Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize).Find(&res.List).Error; err != nil {
		err = s.hm.GetMessage(1001, err)
		return nil, err
	}
	res.Pager = p.NewPager(p.UrlPath)
	return
}

func (s *service) GetArticle(id int64) (article *model.Article, err error) {
	article = new(model.Article)
	article.ID = id
	if err = s.db.Model(article).Preload("User").Preload("Tags").
		Find(article).Error; err != nil {
		err = s.hm.GetMessage(1001, err)
		return nil, err
	}
	var prev, next model.Article
	err = s.db.Model(&prev).Last(&prev, "id < ?", id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == nil {
		article.Prev = &prev
	}
	err = s.db.Model(&next).First(&next, "id > ?", id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == nil {
		article.Next = &next
	}
	return
}

func (s *service) GetLatestArticles(limit int) (articles []*model.Article, err error) {
	key := s.hm.GetCacheKey(4)
	err = s.mc.Get(key, &articles)
	if err != nil {
		if err = s.db.Model(&articles).Order("created_at desc").Limit(limit).Find(&articles).Error; err != nil {
			err = s.hm.GetMessage(1001, err)
			return
		}
		if err = s.mc.Set(key, &articles); err != nil {
			err = s.hm.GetMessage(1002, err)
			return
		}
	}
	return
}

func (s *service) GetHotArticles(limit int) (articles []*model.Article, err error) {
	key := s.hm.GetCacheKey(5)
	err = s.mc.Get(key, &articles)
	if err != nil {
		if err = s.db.Model(&articles).Order("views desc").Limit(limit).Find(&articles).Error; err != nil {
			err = s.hm.GetMessage(1001, err)
			return
		}
		if err = s.mc.Set(key, &articles); err != nil {
			err = s.hm.GetMessage(1002, err)
			return
		}
	}
	return
}
