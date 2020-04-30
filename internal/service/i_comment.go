package service

import (
	"ginana-blog/internal/model"
	"ginana-blog/library/ecode"
)

func (s *service) GetComments(p *model.Pager, objPK int64) (res *model.Comments, err error) {
	res = new(model.Comments)
	query := s.db.Model(&res.List).Where("obj_pk = ?", objPK).Count(&p.AllCount).Where("reply_fk = 0")
	query.Group("user_id").Count(&res.CountUsers)
	query = query.Order("created_at desc").Preload("Children").Preload("User")
	if err = query.Find(&res.List).Error; err != nil {
		err = ecode.Errorf(s.hm.GetError(1001, err))
		return nil, err
	}
	if err = query.Group("user_id").Count(&res.CountUsers).Error; err != nil {
		err = ecode.Errorf(s.hm.GetError(1001, err))
		return nil, err
	}
	res.Pager = p.NewPager(p.UrlPath)
	for _, parent := range res.List {
		for _, child := range parent.Children {
			s.db.Preload("Parent").Preload("User").Find(child)
			if child.Parent != nil {
				s.db.Preload("User").Find(child.Parent)
			} else {
				child.Parent = new(model.Comment)
			}
		}
	}
	return
}

func (s *service) GetLatestComments(limit int) (comments []*model.Comment, err error) {
	key := s.hm.GetCacheKey(4)
	err = s.mc.Get(key, &comments)
	if err != nil {
		if err = s.db.Model(&comments).Order("created_at desc").
			Preload("User").Preload("Article").
			Limit(limit).Find(&comments).Error; err != nil {
			err = ecode.Errorf(s.hm.GetError(1001, err))
			return
		}
		if err = s.mc.Set(key, &comments); err != nil {
			err = ecode.Errorf(s.hm.GetError(1002, err))
			return
		}
	}
	return
}
