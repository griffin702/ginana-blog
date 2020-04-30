package service

import (
	"errors"
	"fmt"
)

func NewHelperMap() (hm HelperMap, err error) {
	hm = &helperMap{
		ErrorHelper: map[int]string{
			500:  "服务器错误",
			1001: "查询失败",
			1002: "创建失败",
			1003: "更新失败",
			1004: "删除失败",
		},
		CacheKey: map[int]string{
			1: "siteOptions",
			2: "latestArticles",
			3: "hotArticles",
			4: "latestComments",
			5: "allLinks",
			6: "user",
			7: "role",
		},
	}
	return
}

type HelperMap interface {
	GetError(i int, args ...string) (int, error)
	GetCacheKey(i int, args ...int64) string
}

type helperMap struct {
	ErrorHelper map[int]string
	CacheKey    map[int]string
}

func (hm *helperMap) GetError(i int, args ...string) (int, error) {
	if len(args) > 1 {
		panic("too many arguments")
	}
	msg := hm.ErrorHelper[i]
	if len(args) == 1 {
		msg = args[0]
	}
	return i, errors.New(msg)
}

func (hm *helperMap) GetCacheKey(i int, args ...int64) string {
	if len(args) > 1 {
		panic("too many arguments")
	}
	key := hm.CacheKey[i]
	if len(args) == 1 {
		key = fmt.Sprintf("%s_%d", key, args[0])
	}
	return key
}
