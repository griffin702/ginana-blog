package service

import (
	"ginana-blog/internal/model"
	"ginana-blog/library/ecode"
)

func (s *service) GetLinks() (res *model.Links, err error) {
	res = new(model.Links)
	query := s.db.Model(&res.List)
	query = query.Order("rank desc")
	if err = query.Find(&res.List).Error; err != nil {
		err = ecode.Errorf(s.GetError(501, err.Error()))
		return nil, err
	}
	return
}
