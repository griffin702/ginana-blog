package service

import (
	"ginana-blog/internal/model"
)

func (s *service) GetComments(p *model.Pager, prs ...model.CommentQueryParam) (res *model.Comments, err error) {
	var pr model.CommentQueryParam
	if len(prs) > 0 {
		pr = prs[0]
	}
	if pr.Order == "" {
		pr.Order = "id desc"
	}
	res = new(model.Comments)
	query := s.db.Model(&res.List)
	if !pr.Admin {
		query = query.Where("obj_pk = ?", pr.ArticleID)
	}
	query.Count(&p.AllCount)
	if !pr.Admin {
		query = query.Where("reply_fk = 0").Preload("Children")
	} else {
		query = query.Preload("Article")
	}
	query.Group("user_id").Count(&res.CountUsers)
	query = query.Order(pr.Order).Preload("User")
	if err = query.Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize).Find(&res.List).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	if !pr.Admin {
		res.LoadParent(s.db)
		p.SetArticleID(pr.ArticleID)
	}
	res.Pager = p
	return
}

func (s *service) GetComment(id int64) (comment *model.Comment, err error) {
	comment = new(model.Comment)
	if err = s.db.Find(comment, "id = ?", id).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	return
}

func (s *service) GetLatestComments(limit int) (comments []*model.Comment, err error) {
	key := s.hm.GetCacheKey(6)
	err = s.mc.Get(key, &comments)
	if err != nil {
		if err = s.db.Model(&comments).Order("id desc").
			Preload("User").Preload("Article").
			Limit(limit).Find(&comments).Error; err != nil {
			return nil, s.hm.GetMessage(1001, err)
		}
		if err = s.mc.Set(key, &comments); err != nil {
			return nil, s.hm.GetMessage(1002, err)
		}
	}
	return
}

func (s *service) CreateComment(req *model.CreateCommentReq) (err error) {
	comment := new(model.Comment)
	comment.ObjPK = req.ObjPK
	comment.ReplyPK = req.ReplyPK
	comment.ReplyFK = req.ReplyFK
	comment.Content = req.Content
	comment.ObjPKType = req.ObjPKType
	comment.IPAddress = req.IPAddress
	comment.UserID = req.UserID
	if err = s.db.Create(comment).Error; err != nil {
		return s.hm.GetMessage(1002, err)
	}
	s.mc.Delete(s.hm.GetCacheKey(6))
	return
}

func (s *service) UpdateComment(req *model.UpdateCommentReq) (err error) {
	comment := new(model.Comment)
	comment.ID = req.ID
	if err = s.db.Find(comment).Error; err != nil {
		return s.hm.GetMessage(1001, err)
	}
	if err = s.db.Model(comment).Update("content", req.Content).Error; err != nil {
		return s.hm.GetMessage(1002, err)
	}
	s.mc.Delete(s.hm.GetCacheKey(6))
	return
}

func (s *service) DeleteComment(id int64) (err error) {
	comment := new(model.Comment)
	if err = s.db.Delete(comment, "id = ?", id).Error; err != nil {
		return s.hm.GetMessage(1004, err)
	}
	s.mc.Delete(s.hm.GetCacheKey(6))
	return
}
