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
	UserID      int64
}

func (c *CAdmin) BeginRequest(ctx iris.Context) {
	tokenStr := c.Session.GetString("token")
	user := c.getUserByToken(tokenStr)
	if user.ID > 0 {
		c.UserID = user.ID
		return
	}
	tokenStr = c.Ctx.GetCookie("token")
	user = c.getUserByToken(tokenStr)
	c.UserID = user.ID
}

func (c *CAdmin) EndRequest(ctx iris.Context) {}

func (c *CAdmin) IsLogin() bool {
	return c.UserID > 0
}

func (c *CAdmin) setHeadMetas(params ...string) {
	c.Ctx.ViewData("isLogin", c.IsLogin())
	title := fmt.Sprintf("inana 后台管理 v%s", c.Config.Version)
	if len(params) > 0 {
		title = fmt.Sprintf("%s - %s", params[0], title)
	}
	c.Ctx.ViewData("title", title)
}

func (c *CAdmin) getUserByToken(tokenStr string) (user *model.UserSession) {
	user = new(model.UserSession)
	token, err := c.Tool.JwtParse(tokenStr, c.Config.JwtSecret)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		c.Session.Set("token", tokenStr)
		if userId, ok := claims["userId"].(float64); ok {
			user.ID = int64(userId)
		}
		if username, ok := claims["username"].(string); ok {
			user.Username = username
			c.Ctx.ViewData("username", username)
		}
	}
	return
}
