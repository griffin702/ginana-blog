package model

import (
	"fmt"
	"math"
	"time"
)

//评论模型
type Comment struct {
	ID        int64      `json:"id" gorm:"primary_key;comment:'评论ID'"`
	CreatedAt time.Time  `json:"created_at" gorm:"comment:'创建时间'"`
	DeletedAt *time.Time `json:"-" sql:"index" gorm:"comment:'删除时间戳'"`
	ObjPK     int64      `json:"obj_pk" form:"obj_pk" gorm:"comment:'article_id'" valid:"gte=0"`
	ReplyPK   int64      `json:"reply_pk" form:"reply_pk" gorm:"index;comment:'comment_id'" valid:"gte=0"`
	ReplyFK   int64      `json:"reply_fk" form:"reply_fk" gorm:"index;comment:'parent_id'" valid:"gte=0"`
	Content   string     `json:"content" form:"content" gorm:"type:LONGTEXT;not null;comment:'评论内容'" valid:"required,max=200"`
	ObjPKType int8       `json:"obj_pk_type" form:"obj_pk_type" gorm:"index;not null;comment:'评论类型'" valid:"gte=0,lte=1"` //0-文章评论，1-友链评论
	IPAddress string     `json:"ip_address" gorm:"type:VARCHAR(255);comment:'IP地址'" valid:"ip"`
	UserID    int64      `json:"user_id" gorm:"index;comment:'user_id'" valid:"required,gt=0"`
	User      *User      `json:"user" gorm:"ForeignKey:UserID"`
	Article   *Article   `json:"article" gorm:"ForeignKey:ObjPK"`
	Parent    *Comment   `json:"parent" gorm:"ForeignKey:ID;AssociationForeignKey:ReplyPK"`
	Children  []*Comment `json:"children" gorm:"ForeignKey:ReplyFK;AssociationForeignKey:ID"`
}

type Comments struct {
	List       []*Comment `json:"list"`
	Pager      *Pager     `json:"pager"`
	CountUsers int64
}

func (m *Comment) ReturnLimit(key string) string {
	data := []rune(m.Content)
	if len(data) > 25 {
		return string(data[:25]) + "..."
	}
	return m.Content
}

//title换行显示
func (m *Comment) Titleln() string {
	data := []rune(m.Content)
	if len(data) > 20 {
		for num := int(math.Floor(float64(len(data)) / 20)); num > 0; num-- {
			copydata := make([]rune, 0)
			copydata = append(copydata, data[:20*num]...)
			copydata = append(copydata, []rune("\n")...)
			data = append(copydata, data[20*num:]...)
		}
		return string(data)
	}
	return m.Content
}

func (m *Comment) ShowSubTime() string {
	sub := int(time.Since(m.CreatedAt).Minutes())
	if sub <= 10 {
		return "刚刚"
	} else if sub > 10 && sub <= 60 {
		return "10分钟前"
	} else if sub > 60 && sub <= 1440 {
		return fmt.Sprintf("%d小时前", sub/60)
	} else if sub > 1440 && sub <= 43200 {
		return fmt.Sprintf("%d天前", sub/1440)
	} else if sub > 43200 && sub <= 525600 {
		return fmt.Sprintf("%d个月前", sub/43200)
	} else {
		return fmt.Sprintf("%d年前", sub/525600)
	}
}
