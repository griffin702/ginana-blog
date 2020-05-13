package model

type Policy struct {
	ID          int64          `json:"id" gorm:"primary_key;comment:'规则ID'"`
	Name        string         `json:"name" gorm:"type:VARCHAR(50);not null;comment:'路由名称'"`
	Router      string         `json:"router" gorm:"type:VARCHAR(191);not null;comment:'请求路由'"`
	Method      string         `json:"method" gorm:"type:VARCHAR(10);not null;comment:'请求方式'"`
	RolePolices []*RolePolices `gorm:"ForeignKey:PolicyID"`
}

type Polices struct {
	List  []*Policy `json:"list"`
	Pager *Pager    `json:"pager"`
}

type CreatePolicyReq struct {
	Name   string `form:"name" valid:"required"`
	Router string `form:"router" valid:"required"`
	Method string `form:"method" valid:"required"`
}

type UpdatePolicyReq struct {
	ID     int64  `form:"id" valid:"required,gt=0"`
	Name   string `form:"name" valid:"required"`
	Router string `form:"router" valid:"required"`
	Method string `form:"method" valid:"required"`
}

type PolicyQueryParam struct {
	Order string
}
