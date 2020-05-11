package service

import (
	"ginana-blog/internal/model"
)

func (s *service) GetMoods(p *model.Pager) (res *model.Moods, err error) {
	res = new(model.Moods)
	query := s.db.Model(&res.List).Count(&p.AllCount)
	query = query.Order("created_at desc")
	if err = query.Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize).Find(&res.List).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	res.Pager = p
	return
}

func (s *service) CreateMood(req *model.MoodReq) (err error) {
	mood := new(model.Mood)
	mood.Content = req.ContentMarkdownDoc
	mood.Cover = req.Cover
	if err = s.db.Create(mood).Error; err != nil {
		return s.hm.GetMessage(1002)
	}
	return
}
