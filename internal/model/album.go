package model

import (
	"strings"
	"time"
)

//相册模型
type Album struct {
	ID        int64     `json:"id" gorm:"primary_key;comment:'相册ID'"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:'创建时间'"`
	Name      string    `json:"name" gorm:"index;not null;type:VARCHAR(100);comment:'相册名称'"`
	Cover     string    `json:"cover" gorm:"type:VARCHAR(255);not null;default:'/static/upload/default/blog-default-0.png';comment:'相册封面'"`
	Hidden    bool      `json:"hidden" gorm:"index;not null;comment:'是否隐藏'"`
	Rank      int8      `json:"rank" gorm:"not null;comment:'权重'"`
	Photos    []*Photo  `json:"photos"`
}

//照片模型
type Photo struct {
	ID        int64     `json:"id" gorm:"primary_key;comment:'照片ID'"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:'创建时间'"`
	AlbumID   int64     `json:"album_id" gorm:"index;not null;comment:'相册ID'"`
	Desc      string    `json:"desc" gorm:"type:VARCHAR(255);not null;comment:'描述'"`
	Url       string    `json:"url" gorm:"type:VARCHAR(255);not null;comment:'URL地址'"`
	Small     string    `json:"small" gorm:"-"`
}

type Albums struct {
	List  []*Album `json:"list"`
	Pager *Pager   `json:"pager"`
}

type Photos struct {
	List  []*Photo `json:"list"`
	Pager *Pager   `json:"pager"`
}

func (m *Album) LongNameAlter() string {
	data := []rune(m.Name)
	length := len(data)
	if length > 15 {
		return string(data[:6]) + "..." + string(data[length-7:length])
	}
	return m.Name
}

func (m *Album) CountPhotos() int {
	return len(m.Photos)
}

func (m *Photo) ChangetoSmall() string {
	arr1 := strings.Split(m.Url, "/")
	filename := arr1[len(arr1)-1]
	arr2 := strings.Split(filename, ".")
	ext := "." + arr2[len(arr2)-1]
	m.Small = strings.Replace(m.Url, ext, "_small"+ext, 1)
	return m.Small
}
