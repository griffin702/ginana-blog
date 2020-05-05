package service

import (
	"ginana-blog/internal/model"
)

func (s *service) GetMoods(p *model.Pager) (res *model.Moods, err error) {
	res = new(model.Moods)
	query := s.db.Model(&res.List).Count(&p.AllCount)
	query = query.Order("created_at desc")
	if err = query.Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize).Find(&res.List).Error; err != nil {
		err = s.hm.GetMessage(1001, err)
		return nil, err
	}
	res.Pager = p.NewPager(p.UrlPath)
	return
}
