package model

import "strings"

type Options struct {
	ID    int64  `json:"id" gorm:"primary_key;comment:'设置ID'"`
	Name  string `json:"name" gorm:"type:VARCHAR(50);unique;not null;comment:'设置名称'"`
	Value string `json:"value" gorm:"type:LONGTEXT;not null;comment:'设置内容'"`
}

type Option struct {
	SiteName     string `form:"site_name"`
	SiteURL      string `form:"site_url"`
	SubTitle     string `form:"sub_title"`
	PageSize     string `form:"page_size"`
	Keywords     string `form:"keywords"`
	Description  string `form:"description"`
	Theme        string `form:"theme"`
	WeiBo        string `form:"wei_bo"`
	Github       string `form:"github"`
	AlbumSize    string `form:"album_size"`
	Nickname     string `form:"nickname"`
	MyOldCity    string `form:"my_old_city"`
	MyCity       string `form:"my_city"`
	MyBirth      string `form:"my_birth"`
	MyProfession string `form:"my_profession"`
	MyLang       string `form:"my_lang"`
	MyLike       string `form:"my_like"`
	MyWorkDesc   string `form:"my_work_desc"`
}

func (o *Option) GetNickname(p ...int) string {
	def := 0
	if len(p) > 0 {
		def = p[0]
	}
	list := strings.Split(o.Nickname, "|")
	return list[def]
}
