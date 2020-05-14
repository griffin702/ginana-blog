package model

import (
	"regexp"
	"strings"
	"time"
)

//心情表
type Mood struct {
	ID        int64     `json:"id" gorm:"primary_key;comment:'心情ID'"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:'创建时间'"`
	Content   string    `json:"content" form:"content" gorm:"type:LONGTEXT;not null;comment:'心情内容'"`
	Cover     string    `json:"cover" form:"cover" gorm:"type:VARCHAR(100);comment:'心情插图'"`
}

type MoodReq struct {
	ContentMarkdownDoc string `form:"mood-content-markdown-doc" valid:"required"`
	ContentHtmlCode    string `form:"mood-content-html-code" valid:"omitempty"`
	Cover              string `form:"cover" valid:"omitempty"`
}

type Moods struct {
	List  []*Mood `json:"list"`
	Pager *Pager  `json:"pager"`
}

func (m *Mood) ChangetoSmall() string {
	arr1 := strings.Split(m.Cover, "/")
	filename := arr1[len(arr1)-1]
	arr2 := strings.Split(filename, ".")
	ext := "." + arr2[len(arr2)-1]
	small := strings.Replace(m.Cover, ext, "_small"+ext, 1)
	return small
}

func (m *Mood) DeleteSmall() string {
	return strings.Replace(m.Cover, "_small", "", 1)
}

func (m *Mood) GetDesc() string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("<[\\S\\s]+?>")
	rep := re.ReplaceAllStringFunc(m.Content, strings.ToLower)
	//去除所有尖括号内的HTML代码
	re, _ = regexp.Compile("<[\\S\\s]+?>")
	rep = re.ReplaceAllString(rep, "")
	return rep
}
