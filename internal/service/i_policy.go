package service

import (
	"ginana-blog/internal/model"
)

func (s *service) GetPolicy(id int64) (policy *model.Policy, err error) {
	policy = new(model.Policy)
	if err = s.db.Preload("RolePolices").Find(policy, "id = ?", id).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	return
}

func (s *service) GetPolices(p *model.Pager, prs ...model.PolicyQueryParam) (res *model.Polices, err error) {
	var pr model.PolicyQueryParam
	if len(prs) > 0 {
		pr = prs[0]
	}
	if pr.Order == "" {
		pr.Order = "id desc"
	}
	res = new(model.Polices)
	query := s.db.Model(&res.List)
	query.Count(&p.AllCount)
	query = query.Order(pr.Order)
	if err = query.Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize).Find(&res.List).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	res.Pager = p
	return
}

func (s *service) GetAllPolices() (res *model.Polices, err error) {
	res = new(model.Polices)
	err = s.db.Find(&res.List).Error
	if err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	return
}

func (s *service) CreatePolicy(req *model.CreatePolicyReq) (policy *model.Policy, err error) {
	policy = new(model.Policy)
	policy.Name = req.Name
	policy.Router = req.Router
	policy.Method = req.Method
	if err = s.db.Create(policy).Error; err != nil {
		return nil, s.hm.GetMessage(1002, err)
	}
	return
}

func (s *service) UpdatePolicy(req *model.UpdatePolicyReq) (policy *model.Policy, err error) {
	policy, err = s.GetPolicy(req.ID)
	if err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	policy.Name = req.Name
	policy.Router = req.Router
	policy.Method = req.Method
	err = s.db.Model(policy).Select("name", "router", "method").Update(policy).Error
	if err != nil {
		return nil, s.hm.GetMessage(1003, err)
	}
	for _, rolePolices := range policy.RolePolices {
		s.mc.Delete(s.hm.GetCacheKey(2, rolePolices.RoleID))
	}
	return
}

func (s *service) DeletePolicy(id int64) (err error) {
	policy, err := s.GetPolicy(id)
	if err != nil {
		return
	}
	tx := s.db.Begin()
	if err = tx.Delete(policy.RolePolices, "policy_id = ?", policy.ID).Error; err != nil {
		tx.Rollback()
		return s.hm.GetMessage(1004, err)
	}
	if err = s.db.Delete(policy, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return s.hm.GetMessage(1004, err)
	}
	tx.Commit()
	for _, rp := range policy.RolePolices {
		s.mc.Delete(s.hm.GetCacheKey(2, rp.RoleID))
	}
	return
}
