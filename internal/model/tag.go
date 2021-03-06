package model

import (
	"fmt"
	"math/rand"
	"time"
)

type Tag struct {
	ID        int64      `json:"id" gorm:"primary_key;comment:'标签ID'"`
	CreatedAt time.Time  `json:"created_at" gorm:"comment:'创建时间'"`
	Name      string     `json:"name" gorm:"type:VARCHAR(50);unique;index;not null;comment:'标签名称'"`
	Articles  []*Article `json:"articles" gorm:"many2many:article_tags"`
}

type TagQueryParam struct {
	Order string
	Admin bool
}

type Tags struct {
	List  []*Tag `json:"list"`
	Pager *Pager `json:"pager"`
}

type TagListReq struct {
	Option  string  `form:"option" valid:"required"`
	NewName string  `form:"new_name" valid:"omitempty"`
	IDs     []int64 `form:"ids" valid:"required,gt=0"`
}

func (t *Tag) Link() string {
	return fmt.Sprintf("<a class=\"category\" href=\"/category/%d\"><span class=\"badge\">%s</span></a>", t.ID, t.Name)
}

func (t *Tag) CountArticles() int {
	return len(t.Articles)
}

//随机一个颜色
func (t *Tag) RangeColor() string {
	var chars = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}
	str := ""
	for i := 0; i < 6; i++ {
		id := rand.Intn(15)
		str += chars[id]
	}
	return "#" + str
}
