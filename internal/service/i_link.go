package service

import (
	"ginana-blog/internal/model"
	"ginana-blog/library/ecode"
)

func (s *service) GetLinks() (links []*model.Link, err error) {
	key := "AllLinks"
	err = s.mc.Get(key, &links)
	if err != nil {
		if err = s.db.Model(&links).Order("created_at desc").Find(&links).Error; err != nil {
			err = ecode.Errorf(s.GetError(1001, err.Error()))
			return
		}
		if err = s.mc.Set(key, &links); err != nil {
			err = ecode.Errorf(s.GetError(1002, err.Error()))
			return
		}
	}
	return
}
