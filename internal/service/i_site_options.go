package service

import (
	"ginana-blog/internal/model"
	"reflect"
)

func (s *service) GetSiteOptions() (res *model.Option, err error) {
	key := s.hm.GetCacheKey(3)
	option := new(model.Option)
	err = s.mc.Get(key, option)
	if err != nil {
		var options []*model.Options
		if err = s.db.Find(&options).Error; err != nil {
			return nil, s.hm.GetMessage(1001, err)
		}
		if len(options) == 0 {
			return nil, s.hm.GetMessage(500, "站点设置为空")
		}
		optionsMap := make(map[string]string)
		for _, v := range options {
			optionsMap[v.Name] = v.Value
		}
		if err = s.tool.MapToStruct(&optionsMap, option); err != nil {
			return nil, s.hm.GetMessage(500, err)
		}
		if err = s.mc.Set(key, option); err != nil {
			return nil, s.hm.GetMessage(1002, err)
		}
	}
	return option, nil
}

func (s *service) UpdateSiteOptions(req *model.Option) (err error) {
	options, err := s.GetSiteOptions()
	if err != nil {
		return s.hm.GetMessage(1001, err)
	}
	ot := reflect.TypeOf(options).Elem()
	ov := reflect.ValueOf(options).Elem()
	rt := reflect.TypeOf(req).Elem()
	rv := reflect.ValueOf(req).Elem()
	tx := s.db.Begin()
	for i := 0; i < ot.NumField(); i++ {
		name := ot.Field(i).Name
		for j := 0; j < rt.NumField(); j++ {
			value := rv.Field(j).Interface()
			if name == rt.Field(j).Name && ov.Field(i).Interface() != value {
				if err = s.db.Model(new(model.Options)).Where("name = ?", name).
					Update("value", value).Error; err != nil {
					tx.Rollback()
					return s.hm.GetMessage(1003, err)
				}
			}
		}
	}
	tx.Commit()
	s.mc.Delete(s.hm.GetCacheKey(3))
	return
}
