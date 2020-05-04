package admin

import (
	"ginana-blog/internal/model"
	"ginana-blog/internal/service"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"strings"
)

type CAdmin struct {
	Ctx       iris.Context
	Session   *sessions.Session
	Svc       service.Service
	Pager     *model.Pager
	GetOption model.GetOption
	Links     *model.Links
	Hm        service.HelperMap
	Valid     model.Validator
}

//获取用户IP地址
func (c *CAdmin) getClientIp() string {
	s := c.Ctx.GetHeader("X-Real-IP")
	if s == "" {
		forwarded := c.Ctx.GetHeader("X-Forwarded-For")
		if forwarded != "" {
			list := strings.Split(forwarded, ":")
			if len(list) > 0 {
				s = list[0]
			}
		} else {
			s = strings.Split(c.Ctx.RemoteAddr(), ":")[0]
		}
	}
	return s
}
