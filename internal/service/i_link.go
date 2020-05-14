package service

import (
	"ginana-blog/internal/model"
)

func (s *service) GetLinks() (links []*model.Link, err error) {
	key := s.hm.GetCacheKey(7)
	err = s.mc.Get(key, &links)
	if err != nil {
		if err = s.db.Model(&links).Order("id desc").Find(&links).Error; err != nil {
			return nil, s.hm.GetMessage(1001, err)
		}
		if err = s.mc.Set(key, &links); err != nil {
			return nil, s.hm.GetMessage(1002, err)
		}
	}
	return
}

func (s *service) GetLink(id int64) (link *model.Link, err error) {
	link = new(model.Link)
	if err = s.db.Find(link, "id = ?", id).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	return
}

func (s *service) CreateLink(req *model.CreateLinkReq) (link *model.Link, err error) {
	link = new(model.Link)
	link.SiteName = req.SiteName
	link.SiteAvatar = req.SiteAvatar
	link.SiteDesc = req.SiteDesc
	link.Url = req.Url
	link.Rank = req.Rank
	if err = s.db.Create(link).Error; err != nil {
		return nil, s.hm.GetMessage(1002, err)
	}
	s.mc.Delete(s.hm.GetCacheKey(7))
	return
}

func (s *service) UpdateLink(req *model.UpdateLinkReq) (link *model.Link, err error) {
	link = new(model.Link)
	link.ID = req.ID
	if err = s.db.Find(link).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	link.SiteName = req.SiteName
	link.SiteAvatar = req.SiteAvatar
	link.SiteDesc = req.SiteDesc
	link.Url = req.Url
	link.Rank = req.Rank
	m, err := s.tool.StructToMap(link)
	if err != nil {
		return nil, s.hm.GetMessage(500, err)
	}
	if err = s.db.Model(link).Update(m).Error; err != nil {
		return nil, s.hm.GetMessage(1003, err)
	}
	s.mc.Delete(s.hm.GetCacheKey(7))
	return
}

func (s *service) DeleteLink(id int64) (err error) {
	link := new(model.Link)
	if err = s.db.Delete(link, "id = ?", id).Error; err != nil {
		return s.hm.GetMessage(1004, err)
	}
	s.mc.Delete(s.hm.GetCacheKey(7))
	return
}
