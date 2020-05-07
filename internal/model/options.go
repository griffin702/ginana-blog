package model

type Options struct {
	ID    int64  `json:"id" gorm:"primary_key;comment:'设置ID'"`
	Name  string `json:"name" gorm:"type:VARCHAR(50);unique;not null;comment:'设置名称'"`
	Value string `json:"value" gorm:"type:LONGTEXT;not null;comment:'设置内容'"`
}

type Option struct {
	SiteName     string `json:"site_name" form:"site_name"`
	SiteURL      string `json:"site_url" form:"site_url"`
	SubTitle     string `json:"sub_title" form:"sub_title"`
	PageSize     string `json:"page_size" form:"page_size"`
	Keywords     string `json:"keywords" form:"keywords"`
	Description  string `json:"description" form:"description"`
	Theme        string `json:"theme" form:"theme"`
	WeiBo        string `json:"wei_bo" form:"wei_bo"`
	Github       string `json:"github" form:"github"`
	AlbumSize    string `json:"album_size" form:"album_size"`
	Nickname     string `json:"nickname" form:"nickname"`
	MyOldCity    string `json:"my_old_city" form:"my_old_city"`
	MyCity       string `json:"my_city" form:"my_city"`
	MyBirth      string `json:"my_birth" form:"my_birth"`
	MyProfession string `json:"my_profession" form:"my_profession"`
	MyLang       string `json:"my_lang" form:"my_lang"`
	MyLike       string `json:"my_like" form:"my_like"`
	MyWorkDesc   string `json:"my_work_desc" form:"my_work_desc"`
}
