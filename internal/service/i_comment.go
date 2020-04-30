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
		err = ecode.Errorf(s.GetError(501, err.Error()))
		return nil, err
	}
	if err = query.Group("user_id").Count(&res.CountUsers).Error; err != nil {
		err = ecode.Errorf(s.GetError(501, err.Error()))
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
