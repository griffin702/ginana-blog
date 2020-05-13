package model

import (
	"time"
)

type Role struct {
	ID        int64     `json:"id" gorm:"primary_key;comment:'角色ID'"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:'创建时间'"`
	RoleName  string    `json:"role_name" gorm:"type:VARCHAR(191);unique;not null;comment:'角色名称'"`
	Polices   []*Policy `json:"polices" gorm:"many2many:role_polices"`
}

type Roles struct {
	List  []*Role `json:"list"`
	Pager *Pager  `json:"pager"`
}

type CreateRoleReq struct {
	RoleName string  `form:"role_name" valid:"required"`
	IDs      []int64 `form:"ids" valid:"omitempty,gt=0"`
}

type UpdateRoleReq struct {
	ID       int64   `form:"id" valid:"required,gt=0"`
	RoleName string  `form:"role_name" valid:"required"`
	IDs      []int64 `form:"ids" valid:"omitempty,gt=0"`
}

type RoleQueryParam struct {
	Order string
}

type RolePolices struct {
	RoleID   int64 `json:"role_id" gorm:"comment:'角色ID'"`
	PolicyID int64 `json:"policy_id" gorm:"comment:'规则ID'"`
}

func (r *Role) CheckPolicy(pid int64) bool {
	for _, policy := range r.Polices {
		if policy.ID == pid {
			return true
		}
	}
	return false
}
