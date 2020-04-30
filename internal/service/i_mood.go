package service

import (
	"ginana-blog/internal/model"
	"ginana-blog/library/ecode"
)

func (s *service) GetMoods(p *model.Pager) (res *model.Moods, err error) {
	res = new(model.Moods)
	query := s.db.Model(&res.List).Count(&p.AllCount)
	query = query.Order("created_at desc")
	if err = query.Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize).Find(&res.List).Error; err != nil {
		err = ecode.Errorf(s.GetError(501, err.Error()))
		return nil, err
	}
	res.Pager = p.NewPager(p.UrlPath)
	return
}
