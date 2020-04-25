package model

import (
	"fmt"
	"time"
)

type Tag struct {
	ID        int64     `json:"id" gorm:"primary_key;comment:'标签ID'"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:'创建时间'"`
	Name      string    `json:"name" gorm:"type:VARCHAR(50);unique;index;not null;comment:'标签名称'"`
}

func (t *Tag) Link() string {
	return fmt.Sprintf("<a class=\"category\" href=\"/category/%d\">%s</a>", t.ID, t.Name)
}
