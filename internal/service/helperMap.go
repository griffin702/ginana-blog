package service

import (
	"errors"
	"fmt"
	"github.com/griffin702/ginana/library/ecode"
)

func NewHelperMap() (hm HelperMap, err error) {
	hm = &helperMap{
		ErrorHelper: map[int]string{
			0:    "请求成功",
			500:  "服务器错误",
			1001: "查询失败",
			1002: "创建失败",
			1003: "更新失败",
			1004: "删除失败",
			1005: "生成验证码失败",
			1006: "验证码失效或发生意外，请重载",
			1007: "验证码校验不正确",
			1008: "密码校验不正确",
			1009: "该帐号未激活",
		},
		CacheKey: map[int]string{
			1: "user",
			2: "role",
			3: "siteOptions",
			4: "latestArticles",
			5: "hotArticles",
			6: "latestComments",
			7: "allLinks",
			8: "captchaCode",
		},
	}
	return
}

type HelperMap interface {
	GetMessage(i int, args ...interface{}) error
	GetCacheKey(i int, args ...interface{}) string
}

type helperMap struct {
	ErrorHelper map[int]string
	CacheKey    map[int]string
}

func (hm *helperMap) GetMessage(i int, args ...interface{}) error {
	msg := hm.ErrorHelper[i]
	var arg interface{}
	if len(args) > 0 {
		arg = args[0]
		if err, ok := arg.(error); ok {
			return ecode.Errorf(i, err)
		}
		if str, ok := arg.(string); ok {
			msg = str
		}
	}
	return ecode.Errorf(i, errors.New(msg))
}

func (hm *helperMap) GetCacheKey(i int, args ...interface{}) string {
	if len(args) > 1 {
		panic("too many arguments")
	}
	key := hm.CacheKey[i]
	if len(args) == 1 {
		key = fmt.Sprintf("%s_%v", key, args[0])
	}
	return key
}
