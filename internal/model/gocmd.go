package model

type PhoneList struct {
	ID        int64  `json:"id" gorm:"primary_key;comment:'用户ID'"`
	ProName   string `json:"pro_name" gorm:"type:VARCHAR(191);unique;not null;comment:'项目名称'"`
	PhoneList string `json:"phone_list" gorm:"type:VARCHAR(255);not null;comment:'手机列表'"`
}

type CreatePhoneListReq struct {
	ProName   string `form:"pro_name" valid:"required"`
	PhoneList string `form:"phone_list" valid:"required"`
}

type UpdatePhoneListReq struct {
	ID        int64  `form:"id" valid:"required,gt=0"`
	ProName   string `form:"pro_name" valid:"required"`
	PhoneList string `form:"phone_list" valid:"required"`
}

type PhoneLists struct {
	List  []*PhoneList `json:"list"`
	Pager *Pager       `json:"pager"`
}
