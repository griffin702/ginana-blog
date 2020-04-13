package service

import (
	"ginana-blog/internal/model"
	"ginana-blog/library/ecode"
)

func (s *service) GetSiteOptions() (res []*model.Options, err error) {
	if err = s.db.Find(&res).Error; err != nil {
		err = ecode.Errorf(500, err)
	}
	return
}
