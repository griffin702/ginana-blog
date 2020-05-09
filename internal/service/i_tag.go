package service

import (
	"ginana-blog/internal/model"
)

func (s *service) GetTags() (res *model.Tags, err error) {
	res = new(model.Tags)
	query := s.db.Model(&res.List)
	query = query.Order("id")
	if err = query.Preload("Articles").Find(&res.List).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	return
}

func (s *service) CountTags() (count int64) {
	tag := new(model.Tag)
	s.db.Model(tag).Count(&count)
	return
}

func (s *service) GetTagByName(name string) (tag *model.Tag, err error) {
	tag = new(model.Tag)
	if err = s.db.Find(tag, "name = ?", name).Error; err != nil {
		return tag, s.hm.GetMessage(1001, err)
	}
	return
}
