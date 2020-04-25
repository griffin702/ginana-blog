package model

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"
)

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
	UserID    int64     `json:"user_id" gorm:"comment:'关联用户ID'"`
	User      *User     `json:"user" gorm:"ForeignKey:UserID"`
	Tags      []*Tag    `json:"tags" gorm:"many2many:article_tags"`
}

type Articles struct {
	List  []*Article `json:"article_list"`
	Pager string     `json:"pager"`
}

//管理员角色关联
type ArticleTags struct {
	ArticleID int64 `json:"article_id"` //文章ID
	TagID     int64 `json:"tag_id"`     //标签ID
}

// 带颜色的标题
func (a *Article) ColorTitle() string {
	if a.Color != "" {
		return fmt.Sprintf("<div style=\"color:%s\">%s</div>", a.Color, a.Title)
	} else {
		return a.Title
	}
}

// URL
func (a *Article) Link() string {
	if a.Urlname != "" {
		if a.Urltype == 1 {
			return fmt.Sprintf("/%s", Rawurlencode(a.Urlname))
		}
		return fmt.Sprintf("/article/%s", Rawurlencode(a.Urlname))
	}
	return fmt.Sprintf("/article/%d", a.ID)
}

// Tags链接
func (a *Article) TagsLink() string {
	if len(a.Tags) == 0 {
		return ""
	}
	var buf bytes.Buffer
	buf.WriteString(`<span><span class="glyphicon glyphicon-tags"></span> [`)
	for k, v := range a.Tags {
		if k > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(v.Link())
	}
	buf.WriteString(`]</span>`)
	return buf.String()
}

//摘要
func (a *Article) Excerpt() string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("<[\\S\\s]+?>")
	rep := re.ReplaceAllStringFunc(a.Content, strings.ToLower)
	//将hr标签转换为特殊标记
	re, _ = regexp.Compile("------------")
	rep = re.ReplaceAllString(rep, "_markdown_hr_")
	//去除所有尖括号内的HTML代码
	re, _ = regexp.Compile("<[\\S\\s]+?>")
	rep = re.ReplaceAllString(rep, "")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{1,}")
	rep = re.ReplaceAllString(rep, "")
	//如果断定截取的断点可能会存在中文字符，则需要转为rune后再截取，否则可能会截成乱码
	data := []rune(rep)
	if i := strings.Index(rep, "_markdown_hr_"); i != -1 {
		return rep[:i] + "..."
	} else if i = -1; len(data) > 58 {
		return string(data[:62]) + "..."
	}
	return rep
}

func (a *Article) DelExcerpt() string {
	if i := strings.Index(a.Content, "------------"); i != -1 {
		x := len("------------")
		return a.Content[i+x:]
	}
	return a.Content
}