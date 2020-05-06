package admin

import (
	"fmt"
	"ginana-blog/internal/config"
	"ginana-blog/internal/model"
	"ginana-blog/internal/service"
	"github.com/griffin702/service/jwt-iris"
	"github.com/griffin702/service/tools"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

type CAdmin struct {
	Ctx         iris.Context
	Session     *sessions.Session
	Svc         service.Service
	JsonPlus    model.JsonPlus
	Pager       *model.Pager
	GetClientIP model.GetClientIP
	GetOption   model.GetOption
	Links       *model.Links
	Hm          service.HelperMap
	Valid       model.Validator
	Tool        *tools.Tool
	Config      *config.Config
}

func (c *CAdmin) IsLogin() bool {
	if _, ok := c.Session.Get("userId").(float64); ok {
		return true
	}
	tokenStr := c.Ctx.GetCookie("token")
	if tokenStr != "" {
		token, err := c.Tool.JwtParse(tokenStr, c.Config.JwtSecret)
		if err != nil {
			return false
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			if userId, ok := claims["userId"].(float64); ok {
				c.Session.Set("userId", userId)
			}
			if username, ok := claims["username"].(string); ok {
				c.Session.Set("username", username)
			}
		}
		return true
	}
	return false
}

func (c *CAdmin) setHeadMetas(params ...string) {
	isLogin := c.IsLogin()
	c.Ctx.ViewData("isLogin", isLogin)
	if isLogin {
		title := fmt.Sprintf("inana 后台管理 v%s", c.Config.Version)
		if len(params) > 0 {
			title = fmt.Sprintf("%s - %s", params[0], title)
		}
		c.Ctx.ViewData("title", title)
		c.Ctx.ViewData("userId", c.Session.GetInt64Default("userId", 0))
		c.Ctx.ViewData("username", c.Session.GetStringDefault("username", ""))
	}
}
