package service

import (
	"ginana-blog/internal/model"
)

func (s *service) GetSiteOptions() (res map[string]string, err error) {
	key := s.hm.GetCacheKey(3)
	var options []*model.Options
	err = s.mc.Get(key, &options)
	if err != nil {
		if err = s.db.Find(&options).Error; err != nil {
			err = s.hm.GetMessage(1001, err)
			return
		}
		if len(options) == 0 {
			err = s.hm.GetMessage(500, "站点设置异常")
			return
		}
		if err = s.mc.Set(key, &options); err != nil {
			err = s.hm.GetMessage(1002, err)
			return
		}
	}
	res = make(map[string]string)
	for _, v := range options {
		res[v.Name] = v.Value
	}
	return
}
