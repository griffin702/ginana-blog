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
	key := reflect.TypeOf(options).Elem()
	value := reflect.ValueOf(options).Elem()
	reqKey := reflect.TypeOf(req).Elem()
	reqValue := reflect.ValueOf(req).Elem()
	tx := s.db.Begin()
	for i := 0; i < key.NumField(); i++ {
		name := key.Field(i).Name
		for j := 0; j < reqKey.NumField(); j++ {
			v := reqValue.Field(j).Interface()
			if name == reqKey.Field(j).Name && value.Field(i).Interface() != v {
				opt := new(model.Options)
				err = s.db.Model(opt).Where("name = ?", name).Update("value", v).Error
				if err != nil {
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
