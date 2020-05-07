package service

import (
	"ginana-blog/internal/model"
)

func (s *service) GetSiteOptions() (res *model.Option, err error) {
	key := s.hm.GetCacheKey(3)
	option := new(model.Option)
	err = s.mc.Get(key, option)
	if err != nil {
		var options []*model.Options
		if err = s.db.Find(&options).Error; err != nil {
			err = s.hm.GetMessage(1001, err)
			return
		}
		if len(options) == 0 {
			err = s.hm.GetMessage(500, "站点设置为空")
			return
		}
		optionsMap := make(map[string]string)
		for _, v := range options {
			optionsMap[v.Name] = v.Value
		}
		if err = s.tool.MapToStruct(&optionsMap, option); err != nil {
			err = s.hm.GetMessage(500, err)
			return
		}
		if err = s.mc.Set(key, option); err != nil {
			err = s.hm.GetMessage(1002, err)
			return
		}
	}
	res = option
	return
}

func (s *service) UpdateSiteOptions(req *model.Option) (err error) {
	m, err := s.tool.StructToMap(req)
	if err != nil {
		err = s.hm.GetMessage(1010, err)
		return
	}
	s.db.Begin()
	for k, v := range m {
		if value, ok := v.(string); ok {
			options := new(model.Options)
			if err = s.db.Find(options, "name = ?", k).Error; err != nil {
				err = nil
				continue
			}
			options.Value = value
			if err = s.db.Model(options).Update(options).Error; err != nil {
				err = s.hm.GetMessage(1003, err)
				s.db.Rollback()
				return
			}
		}
	}
	s.db.Commit()
	s.mc.Delete(s.hm.GetCacheKey(3))
	return
}
