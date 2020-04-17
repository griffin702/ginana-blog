package service

import (
	"context"
	"fmt"
	"ginana-blog/internal/model"
	"ginana-blog/library/database"
	"ginana-blog/library/ecode"
	"sync"
)

func (s *service) GetEFUsers(c context.Context) (users []*database.CasbinUser, err error) {
	var userIdList []int64
	s.db.Model(&model.User{}).Select("id").Pluck("id", &userIdList)
	var wg sync.WaitGroup
	var ch = make(chan int64, 1)
	wg.Add(len(userIdList))
	for _, userId := range userIdList {
		go func(userId int64, users *[]*database.CasbinUser, wg *sync.WaitGroup) {
			u, err := s.GetUser(c, userId)
			if err != nil {
				return
			}
			user := new(database.CasbinUser)
			user.ID = u.ID
			for _, role := range u.Roles {
				user.RoleNames = append(user.RoleNames, role.RoleName)
			}
			ch <- userId
			*users = append(*users, user)
			<-ch
			wg.Done()
		}(userId, &users, &wg)
	}
	wg.Wait()
	return
}

func (s *service) GetUser(ctx context.Context, id int64) (user *model.User, err error) {
	key := fmt.Sprintf("user_%d", id)
	user = new(model.User)
	err = s.mc.Get(key, user)
	if err != nil {
		user.ID = id
		if err = s.db.Find(user).Related(&user.Roles, "Roles").Error; err != nil {
			err = ecode.Errorf(s.GetError(1001, err.Error()))
			return
		}
		if err = s.mc.Set(key, user); err != nil {
			err = ecode.Errorf(s.GetError(1002, err.Error()))
			return
		}
	}
	return
}
