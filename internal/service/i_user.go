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

func (s *service) GetUser(id int64) (user *model.User, err error) {
	key := s.hm.GetCacheKey(1, id)
	user = new(model.User)
	err = s.mc.Get(key, user)
	if err != nil {
		user.ID = id
		if err = s.db.Find(user).Related(&user.Roles, "Roles").Error; err != nil {
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
