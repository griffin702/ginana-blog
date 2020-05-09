package model

import (
	"time"
)

type User struct {
	ID          int64      `json:"id" gorm:"primary_key;comment:'用户ID'"`
	CreatedAt   time.Time  `json:"created_at" gorm:"comment:'创建时间'"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"comment:'更新时间'"`
	DeletedAt   *time.Time `json:"-" sql:"index" gorm:"comment:'删除时间戳'"`
	Username    string     `json:"username" gorm:"type:VARCHAR(191);unique;not null;comment:'用户账号'"`
	Password    string     `json:"-" gorm:"type:VARCHAR(255);not null;comment:'用户密码'"`
	Nickname    string     `json:"nickname" gorm:"type:VARCHAR(100);unique;not null;comment:'用户昵称'"`
	Email       string     `json:"email" gorm:"type:VARCHAR(255);not null;comment:'用户邮箱'"`
	Avatar      string     `json:"avatar" gorm:"type:VARCHAR(255);not null;default:'/static/upload/default/user-default-60x60.png';comment:'用户头像'"`
	IsAuth      bool       `json:"is_auth" gorm:"comment:'认证(0-正常,1-未认证)'"`
	LastLoginIP string     `json:"last_login_ip" gorm:"type:VARCHAR(30);not null;comment:'最后登录IP'"`
	CountLogin  int64      `json:"count_login" gorm:"comment:'登录次数'"`
	Roles       []*Role    `json:"roles" gorm:"many2many:user_roles"`
}

type CreateUserReq struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Nickname string `form:"nickname" binding:"omitempty"`
	Email    string `form:"email" binding:"omitempty,email"`
	Avatar   string `form:"avatar" binding:"omitempty"`
	IsAuth   bool   `form:"is_auth"`
}

type UpdateUserReq struct {
	ID       int64  `form:"id" binding:"required,gt=0"`
	Password string `form:"password" binding:"omitempty"`
	Nickname string `form:"nickname" binding:"omitempty"`
	Email    string `form:"email" binding:"omitempty,email"`
	Avatar   string `form:"avatar" binding:"omitempty"`
	IsAuth   bool   `form:"is_auth"`
}

type Users struct {
	List  []*User `json:"list"`
	Pager *Pager  `json:"pager"`
}

type UserQueryParam struct {
	Order string
}

type UserRoles struct {
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
}

type UserLoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Captcha  string `json:"captcha" binding:"required"`
	LoginIP  string `json:"login_ip"`
}

type UserSession struct {
	ID       int64
	Username string
	Roles    []string
}
