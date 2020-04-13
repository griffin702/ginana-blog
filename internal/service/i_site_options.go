package service

import (
	"ginana-blog/internal/model"
	"ginana-blog/library/ecode"
)

func (s *service) GetSiteOptions() (res map[string]*string, err error) {
	var options []*model.Options
	if err = s.db.Find(&options).Error; err != nil {
		err = ecode.Errorf(500, err)
		return
	}
	res = make(map[string]*string)
	for _, v := range options {
		res[v.Name] = &v.Value
	}
	return
}
