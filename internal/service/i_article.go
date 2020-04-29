package service

import (
	"ginana-blog/internal/model"
	"ginana-blog/library/ecode"
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
	query := s.db.Model(&res.List).Count(&p.AllCount)
	query = query.Order(pr.Order)
	if err := query.Preload("User").Preload("Tags").
		Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize).Find(&res.List).Error; err != nil {
		err = ecode.Errorf(s.GetError(501, err.Error()))
		return nil, err
	}
	res.Pager = p.NewPager(p.UrlPath).ToString()
	return
}
