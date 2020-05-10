package service

import (
	"context"
	"ginana-blog/internal/model"
	"github.com/griffin702/ginana/library/database"
	"sync"
)

func (s *service) GetEFUsers(c context.Context) (users []*database.EFUseRole, err error) {
	var userIdList []int64
	s.db.Model(&model.User{}).Select("id").Pluck("id", &userIdList)
	var wg sync.WaitGroup
	var ch = make(chan int64, 1)
	wg.Add(len(userIdList))
	for _, userId := range userIdList {
		go func(userId int64, users *[]*database.EFUseRole, wg *sync.WaitGroup) {
			u, err := s.GetUser(userId)
			if err != nil {
				return
			}
			ch <- userId
			for _, role := range u.Roles {
				user := new(database.EFUseRole)
				user.UserID = u.ID
				user.RoleName = role.RoleName
				*users = append(*users, user)
			}
			<-ch
			wg.Done()
		}(userId, &users, &wg)
	}
	wg.Wait()
	return
}

func (s *service) GetUsers(p *model.Pager, prs ...model.UserQueryParam) (res *model.Users, err error) {
	var pr model.UserQueryParam
	if len(prs) > 0 {
		pr = prs[0]
	}
	if pr.Order == "" {
		pr.Order = "id desc"
	}
	res = new(model.Users)
	query := s.db.Model(&res.List)
	query.Count(&p.AllCount)
	query = query.Order(pr.Order).Preload("Roles")
	if err = query.Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize).Find(&res.List).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	res.Pager = p
	return
}

func (s *service) GetUser(id int64) (user *model.User, err error) {
	key := s.hm.GetCacheKey(1, id)
	user = new(model.User)
	err = s.mc.Get(key, user)
	if err != nil {
		if err = s.db.Find(user, "id = ?", id).Related(&user.Roles, "Roles").Error; err != nil {
			return nil, s.hm.GetMessage(1001, err)
		}
		if err = s.mc.Set(key, user); err != nil {
			return nil, s.hm.GetMessage(1002, err)
		}
	}
	return
}

func (s *service) GetUserByUsername(username string) (user *model.User, err error) {
	user = new(model.User)
	if err = s.db.Find(user, "username = ?", username).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	return
}

func (s *service) CountUsers() (count int64) {
	user := new(model.User)
	s.db.Model(user).Count(&count)
	return
}

func (s *service) CreateUser(req *model.CreateUserReq) (user *model.User, err error) {
	user = new(model.User)
	user.Username = req.Username
	user.Password = s.tool.BcryptHashGenerate(req.Password)
	user.Nickname = req.Nickname
	user.Email = req.Email
	user.Avatar = req.Avatar
	user.IsAuth = req.IsAuth
	if err = s.db.Create(user).Error; err != nil {
		return nil, s.hm.GetMessage(1002, err)
	}
	return
}

func (s *service) UpdateUser(req *model.UpdateUserReq) (user *model.User, err error) {
	user = new(model.User)
	user.ID = req.ID
	if err = s.db.Find(user).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	if req.Password != "" {
		user.Password = s.tool.BcryptHashGenerate(req.Password)
	}
	user.Nickname = req.Nickname
	user.Avatar = req.Avatar
	user.Email = req.Email
	user.IsAuth = req.IsAuth
	m, err := s.tool.StructToMap(user)
	if err != nil {
		return nil, s.hm.GetMessage(500, err)
	}
	if err = s.db.Model(user).Update(m).Error; err != nil {
		return nil, s.hm.GetMessage(1002, err)
	}
	s.mc.Delete(s.hm.GetCacheKey(1, user.ID))
	return
}

func (s *service) UpdateAccount(req *model.UpdateUserReq) (user *model.User, err error) {
	user = new(model.User)
	user.ID = req.ID
	if err = s.db.Find(user).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	if req.Password != "" {
		if !s.tool.BcryptHashCompare(user.Password, req.Password) {
			return nil, s.hm.GetMessage(1008)
		}
		user.Password = s.tool.BcryptHashGenerate(req.NewPassword)
	}
	user.Nickname = req.Nickname
	user.Avatar = req.Avatar
	user.Email = req.Email
	m, err := s.tool.StructToMap(user)
	if err != nil {
		return nil, s.hm.GetMessage(500, err)
	}
	if err = s.db.Model(user).Update(m).Error; err != nil {
		return nil, s.hm.GetMessage(1002, err)
	}
	s.mc.Delete(s.hm.GetCacheKey(1, user.ID))
	return
}

func (s *service) DeleteUser(id int64) (err error) {
	user := new(model.User)
	if err = s.db.Delete(user, "id = ?", id).Error; err != nil {
		return s.hm.GetMessage(1004, err)
	}
	s.mc.Delete(s.hm.GetCacheKey(1, id))
	return
}
