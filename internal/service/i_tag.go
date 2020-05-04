package service

import (
	"ginana-blog/internal/model"
	"github.com/griffin702/ginana/library/ecode"
)

func (s *service) GetTags() (res *model.Tags, err error) {
	res = new(model.Tags)
	query := s.db.Model(&res.List)
	query = query.Order("id")
	if err = query.Preload("Articles").Find(&res.List).Error; err != nil {
		err = ecode.Errorf(s.hm.GetError(1001, err))
		return nil, err
	}
	return
}
