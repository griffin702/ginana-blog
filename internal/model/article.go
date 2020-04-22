package model

import "time"

type Article struct {
	ID        int64     `json:"id" gorm:"primary_key;comment:'文章ID'"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:'创建时间'"`
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:'更新时间'"`
	Title     string    `json:"title" gorm:"type:VARCHAR(255);unique;not null;index;comment:'文章标题'"`
	Color     string    `json:"color" gorm:"type:VARCHAR(10);not null;comment:'标题颜色'"`
	Urlname   string    `json:"urlname" gorm:"type:VARCHAR(100);not null;comment:'特殊链接名称'"`
	Urltype   int8      `json:"urltype" gorm:"comment:'特殊链接类型'"`
	Content   string    `json:"content" gorm:"type:LONGTEXT;not null;comment:'文章内容'"`
	Views     int64     `json:"views" gorm:"comment:'查看次数'"`
	Status    int8      `json:"status" gorm:"comment:'文章状态'"`
	Istop     int8      `json:"istop" gorm:"comment:'置顶相关'"`
	Cover     string    `json:"cover" gorm:"type:VARCHAR(255);default:'/static/upload/default/blog-default-0.png';not null;comment:'文章封面'"`
	User      *User     `json:"user" gorm:"ForeignKey:ID"`
	Tags      []*Tag    `json:"tags" gorm:"many2many:article_tags"`
}

//管理员角色关联
type ArticleTags struct {
	ArticleID int64 `json:"article_id"` //文章ID
	TagID     int64 `json:"tag_id"`     //标签ID
}
