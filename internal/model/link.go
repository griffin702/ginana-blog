package model

import "time"

type Link struct {
	ID         int64     `json:"id" gorm:"primary_key;comment:'友链ID'"`
	CreatedAt  time.Time `json:"created_at" gorm:"comment:'创建时间'"`
	SiteName   string    `json:"site_name" gorm:"type:VARCHAR(80);not null;index;comment:'站点名称'"`
	SiteAvatar string    `json:"site_avatar" gorm:"type:VARCHAR(200);default:'/static/upload/default/user-default-60x60.png';not null;comment:'站点头像'"`
	Url        string    `json:"url" gorm:"type:VARCHAR(200);not null;comment:'站点链接'"`
	SiteDesc   string    `json:"site_desc" gorm:"type:VARCHAR(300);not null;comment:'站点描述'"`
	Rank       int8      `json:"rank" gorm:"not null;comment:'站点权重'"`
}

type CreateLinkReq struct {
	SiteName   string `form:"site_name" valid:"required"`
	SiteAvatar string `form:"site_avatar" valid:"omitempty"`
	Url        string `form:"url" valid:"required"`
	SiteDesc   string `form:"site_desc" valid:"required"`
	Rank       int8   `form:"rank" valid:"omitempty"`
}

type UpdateLinkReq struct {
	ID         int64  `form:"id" valid:"required,gt=0"`
	SiteName   string `form:"site_name" valid:"required"`
	SiteAvatar string `form:"site_avatar" valid:"omitempty"`
	Url        string `form:"url" valid:"required"`
	SiteDesc   string `form:"site_desc" valid:"required"`
	Rank       int8   `form:"rank" valid:"omitempty"`
}

type Links struct {
	List []*Link `json:"list"`
}
