package service

import (
	"context"
	"ginana-blog/internal/model"
	"github.com/griffin702/ginana/library/database"
	"github.com/jinzhu/gorm"
	"sync"
)

func (s *service) GetEFRoles(c context.Context) (roles []*database.EFRolePolicy, err error) {
	var roleIdList []int64
	s.db.Model(&model.Role{}).Select("id").Pluck("id", &roleIdList)
	var wg sync.WaitGroup
	var ch = make(chan int64, 1)
	wg.Add(len(roleIdList))
	for _, roleId := range roleIdList {
		go func(roleId int64, roles *[]*database.EFRolePolicy, wg *sync.WaitGroup) {
			r, err := s.GetRole(roleId)
			if err != nil {
				return
			}
			ch <- roleId
			for _, policy := range r.Polices {
				role := new(database.EFRolePolicy)
				role.RoleName = r.RoleName
				role.Router = policy.Router
				role.Method = policy.Method
				*roles = append(*roles, role)
			}
			<-ch
			wg.Done()
		}(roleId, &roles, &wg)
	}
	wg.Wait()
	return
}

func (s *service) GetRole(id int64) (role *model.Role, err error) {
	key := s.hm.GetCacheKey(2, id)
	role = new(model.Role)
	err = s.mc.Get(key, role)
	if err != nil {
		role.ID = id
		if err = s.db.Find(role).Related(&role.Polices, "Polices").Error; err != nil {
			return nil, s.hm.GetMessage(1001, err)
		}
		if err = s.mc.Set(key, role); err != nil {
			return nil, s.hm.GetMessage(1002, err)
		}
	}
	return
}

func (s *service) GetRoles(p *model.Pager, prs ...model.RoleQueryParam) (res *model.Roles, err error) {
	var pr model.RoleQueryParam
	if len(prs) > 0 {
		pr = prs[0]
	}
	if pr.Order == "" {
		pr.Order = "id desc"
	}
	res = new(model.Roles)
	query := s.db.Model(&res.List)
	query.Count(&p.AllCount)
	query = query.Order(pr.Order).Preload("Polices")
	if err = query.Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize).Find(&res.List).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	res.Pager = p
	return
}

func (s *service) GetRoleByName(name string) (role *model.Role, err error) {
	role = new(model.Role)
	if err = s.db.Find(role, "role_name = ?", name).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	return
}

func (s *service) CreateRole(req *model.CreateRoleReq) (role *model.Role, err error) {
	role, err = s.GetRoleByName(req.RoleName)
	if err == gorm.ErrRecordNotFound {
		role = new(model.Role)
		role.RoleName = req.RoleName
		for _, pid := range req.IDs {
			policy, err := s.GetPolicy(pid)
			if err == gorm.ErrRecordNotFound {
				continue
			} else if err != nil {
				return nil, s.hm.GetMessage(1001, err)
			}
			role.Polices = append(role.Polices, policy)
		}
		if err = s.db.Create(role).Error; err != nil {
			return nil, s.hm.GetMessage(1002, err)
		}
		return
	}
	return nil, s.hm.GetMessage(1001, err)
}

func (s *service) UpdateRole(req *model.UpdateRoleReq) (role *model.Role, err error) {
	role, err = s.GetRole(req.ID)
	if err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	tx := s.db.Begin()
	err = tx.Model(role).Update("role_name = ?", req.RoleName).Error
	if err != nil {
		tx.Rollback()
		return nil, s.hm.GetMessage(1003, err)
	}
	err = tx.Delete(&model.RolePolices{}, "role_id = ? and policy_id in (?)", role.ID, req.IDs).Error
	if err != nil {
		tx.Rollback()
		return nil, s.hm.GetMessage(1004, err)
	}
	for _, pid := range req.IDs {
		rolePolices := new(model.RolePolices)
		rolePolices.RoleID = role.ID
		rolePolices.PolicyID = pid
		if err = tx.Create(rolePolices).Error; err != nil {
			tx.Rollback()
			return nil, s.hm.GetMessage(1002, err)
		}
	}
	tx.Commit()
	s.mc.Delete(s.hm.GetCacheKey(2, role.ID))
	return
}

func (s *service) DeleteRole(id int64) (err error) {
	role, err := s.GetRole(id)
	if err != nil {
		return
	}
	tx := s.db.Begin()
	for _, policy := range role.Polices {
		err = tx.Delete(&model.RolePolices{}, "role_id = ? and policy_id in (?)", role.ID, policy.ID).Error
		if err != nil {
			tx.Rollback()
			return s.hm.GetMessage(1004, err)
		}
	}
	if err = tx.Delete(role, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return s.hm.GetMessage(1004, err)
	}
	tx.Commit()
	s.mc.Delete(s.hm.GetCacheKey(2, role.ID))
	return
}

func (s *service) GetPolicy(id int64) (policy *model.Policy, err error) {
	policy = new(model.Policy)
	if err = s.db.Preload("RolePolices").Find(policy, "id = ?", id).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	return
}

func (s *service) GetPolices() (polices *model.Polices, err error) {
	polices = new(model.Polices)
	if err = s.db.Find(&polices.List).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return
		}
		return nil, s.hm.GetMessage(1001, err)
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
