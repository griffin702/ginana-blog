package service

import (
	"ginana-blog/internal/model"
	"ginana-blog/library/ecode"
)

func (s *service) GetArticles(p *model.Pager) (res *model.Articles, err error) {
	res = new(model.Articles)
	query := s.db.Model(&res.List).Count(&p.AllCount)
	if err := query.Order("id").
		Preload("User").Preload("Tags").
		Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize).Find(&res.List).Error; err != nil {
		err = ecode.Errorf(s.GetError(501, err.Error()))
		return nil, err
	}
	res.Pager = p.NewPager("/").ToString()
	return
}
