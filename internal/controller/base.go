package controller

import (
	"ginana-blog/internal/config"
	"ginana-blog/internal/model"
	"ginana-blog/internal/service"
	"github.com/griffin702/service/jwt-iris"
	"github.com/griffin702/service/tools"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"strings"
)

type BaseController struct {
	Ctx          iris.Context
	Session      *sessions.Session
	Svc          service.Service
	JsonPlus     model.JsonPlus
	Pager        *model.Pager
	GetClientIP  model.GetClientIP
	SiteOptions  *model.Option
	Links        *model.Links
	DisableRight bool
	Hm           service.HelperMap
	Valid        model.Validator
	Tool         *tools.Tool
	Config       *config.Config
	UserID       int64
}

func (c *BaseController) BeginRequest(ctx iris.Context) {
	user := c.GetUserByToken()
	c.UserID = user.ID
}

func (c *BaseController) EndRequest(ctx iris.Context) {}

func (c *BaseController) GetUserByToken() (user *model.UserSession) {
	user = new(model.UserSession)
	tokenStr := c.Session.GetString("token")
	if tokenStr == "" {
		tokenStr = c.Ctx.GetCookie("token")
		if tokenStr == "" {
			return
		}
	}
	token, err := c.Tool.JwtParse(tokenStr, c.Config.JwtSecret)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		c.Session.Set("token", tokenStr)
		if userId, ok := claims["userId"].(float64); ok {
			user.ID = int64(userId)
			if user.ID > 0 {
				c.Ctx.ViewData("isLogin", true)
			}
		}
		if username, ok := claims["username"].(string); ok {
			user.Username = username
			c.Ctx.ViewData("username", username)
		}
	}
	return
}

func (c *BaseController) IsLogin() bool {
	return c.UserID > 0
}

func (c *BaseController) ShowMsg(msg string) {
	redirect := c.Ctx.GetReferrer().Path
	if redirect == "" {
		redirect = "/"
	}
	c.Ctx.ViewData("redirect", redirect)
	c.Ctx.ViewData("message", msg)
	c.Ctx.View("message/message.html")
	c.Ctx.StopExecution()
}

func (c *BaseController) IsDefaultSrc(value string) bool {
	var defaultDir = "/static/upload/default/"
	if value != "" {
		if index := strings.Index(value, defaultDir); index != -1 {
			return true
		}
	}
	return false
}
