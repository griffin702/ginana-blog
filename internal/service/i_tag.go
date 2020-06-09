package service

import (
	"fmt"
	"ginana-blog/internal/model"
	"github.com/jinzhu/gorm"
)

func (s *service) GetTags(p *model.Pager, prs ...model.TagQueryParam) (res *model.Tags, err error) {
	var pr model.TagQueryParam
	if len(prs) > 0 {
		pr = prs[0]
	}
	if pr.Order == "" {
		pr.Order = "id desc"
	}
	res = new(model.Tags)
	query := s.db.Model(&res.List)
	an := s.db.NewScope(&res.List).TableName()
	bn := s.db.NewScope(&model.ArticleTags{}).TableName()
	joinStr := fmt.Sprintf("left join %s on %s.id = %s.tag_id", bn, an, bn)
	groupStr := fmt.Sprintf("%s.id", an)
	orderStr := fmt.Sprintf("count(%s.id) desc", an)
	query = query.Joins(joinStr).Group(groupStr).Order(orderStr).Order(pr.Order)
	if !pr.Admin {
		query = query.Having(fmt.Sprintf("count(%s.id) > 9", an))
	}
	query.Count(&p.AllCount)
	query = query.Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize)
	if err = query.Preload("Articles").Find(&res.List).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	res.Pager = p
	return
}

func (s *service) CountTags() (count int64) {
	tag := new(model.Tag)
	s.db.Model(tag).Count(&count)
	return
}

func (s *service) GetTagByName(name string) (tag *model.Tag, err error) {
	tag = new(model.Tag)
	if err = s.db.Find(tag, "name = ?", name).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	return
}

func (s *service) BatchTag(req *model.TagListReq) (err error) {
	switch req.Option {
	case "merge":
		if req.NewName == "" {
			return s.hm.GetMessage(500, "未输入新标签名")
		}
		var tags []*model.Tag
		err = s.db.Preload("Articles").Find(&tags, "id in (?)", req.IDs).Error
		if err != nil {
			return s.hm.GetMessage(500, err)
		}
		tx := s.db.Begin()
		tag, err := s.GetTagByName(req.NewName)
		if err == gorm.ErrRecordNotFound {
			tag = new(model.Tag)
			tag.Name = req.NewName
			if err = tx.Create(tag).Error; err != nil {
				return s.hm.GetMessage(500, err)
			}
		} else if err != nil {
			return s.hm.GetMessage(500, err)
		} else {
			for i := 0; i < len(req.IDs); i++ {
				if req.IDs[i] == tag.ID {
					req.IDs = append(req.IDs[:i], req.IDs[i+1:]...)
					i--
				}
			}
		}
		if err = tx.Delete(&model.Tag{}, "id in (?)", req.IDs).Error; err != nil {
			tx.Rollback()
			return s.hm.GetMessage(1004, err)
		}
		articles := make(map[int64]bool)
		for _, t := range tags {
			for _, a := range t.Articles {
				articles[a.ID] = true
				if t.ID != tag.ID {
					articles[a.ID] = false
				}
			}
		}
		for id, ok := range articles {
			articleTags := new(model.ArticleTags)
			err = tx.Delete(articleTags, "article_id = ? and tag_id in (?)", id, req.IDs).Error
			if err != nil {
				tx.Rollback()
				return s.hm.GetMessage(1004, err)
			}
			if !ok {
				articleTags.ArticleID = id
				articleTags.TagID = tag.ID
				if err = tx.Create(articleTags).Error; err != nil {
					tx.Rollback()
					return s.hm.GetMessage(1002, err)
				}
			}
		}
		tx.Commit()
	case "delete":
		tag := new(model.Tag)
		if err = s.db.Model(tag).Where("id in (?)", req.IDs).Delete(tag).Error; err != nil {
			return s.hm.GetMessage(1004, err)
		}
	}
	return
}

func (s *service) GetTagsLimit6() (tags []*model.Tag, err error) {
	query := s.db.Model(&tags)
	an := s.db.NewScope(&tags).TableName()
	bn := s.db.NewScope(&model.ArticleTags{}).TableName()
	joinStr := fmt.Sprintf("left join %s on %s.id = %s.tag_id", bn, an, bn)
	groupStr := fmt.Sprintf("%s.id", an)
	orderStr := fmt.Sprintf("count(%s.id) desc", an)
	query = query.Joins(joinStr).Group(groupStr).Order(orderStr).Limit(6)
	if err = query.Preload("Articles", func(db *gorm.DB) *gorm.DB {
		return db.Order("istop desc, id desc")
	}).Find(&tags).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	for _, tag := range tags {
		if len(tag.Articles) > 8 {
			tag.Articles = tag.Articles[:8]
		}
	}
	return
}
