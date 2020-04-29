package service

import (
	"ginana-blog/internal/model"
	"ginana-blog/library/ecode"
)

func (s *service) GetTags() (res *model.Tags, err error) {
	res = new(model.Tags)
	query := s.db.Model(&res.List)
	query = query.Order("id")
	if err := query.Preload("Articles").Find(&res.List).Error; err != nil {
		err = ecode.Errorf(s.GetError(501, err.Error()))
		return nil, err
	}
	return
}
