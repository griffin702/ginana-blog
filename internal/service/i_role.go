package service

import (
	"context"
	"ginana-blog/internal/model"
	"ginana-blog/library/database"
	"ginana-blog/library/ecode"
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
			r, err := s.GetRole(c, roleId)
			if err != nil {
				return
			}
			ch <- roleId
			for _, policy := range r.Policys {
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

func (s *service) GetRole(ctx context.Context, id int64) (role *model.Role, err error) {
	key := s.hm.GetCacheKey(7, id)
	role = new(model.Role)
	err = s.mc.Get(key, role)
	if err != nil {
		role.ID = id
		if err = s.db.Find(role).Related(&role.Policys, "Policys").Error; err != nil {
			err = ecode.Errorf(s.hm.GetError(1001, err))
			return
		}
		if err = s.mc.Set(key, role); err != nil {
			err = ecode.Errorf(s.hm.GetError(1002, err))
			return
		}
	}
	return
}
