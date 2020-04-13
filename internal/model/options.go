package model

type Options struct {
	ID    int64  `json:"id" gorm:"primary_key;comment:'设置ID'"`
	Name  string `json:"name" gorm:"type:VARCHAR(50);unique;not null;comment:'设置名称'"`
	Value string `json:"value" gorm:"type:LONGTEXT;not null;comment:'设置内容'"`
}
