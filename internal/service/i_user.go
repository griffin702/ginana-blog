package service

import (
	"context"
	"ginana-blog/internal/model"
	"ginana-blog/library/database"
	"ginana-blog/library/ecode"
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
			u, err := s.GetUser(c, userId)
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

func (s *service) GetUser(ctx context.Context, id int64) (user *model.User, err error) {
	key := s.hm.GetCacheKey(6, id)
	user = new(model.User)
	err = s.mc.Get(key, user)
	if err != nil {
		user.ID = id
		if err = s.db.Find(user).Related(&user.Roles, "Roles").Error; err != nil {
			err = ecode.Errorf(s.hm.GetError(1001, err.Error()))
			return
		}
		if err = s.mc.Set(key, user); err != nil {
			err = ecode.Errorf(s.hm.GetError(1002, err.Error()))
			return
		}
	}
	return
}
