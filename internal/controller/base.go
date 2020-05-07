package controller

import (
	"ginana-blog/internal/config"
	"ginana-blog/internal/model"
	"ginana-blog/internal/service"
	"github.com/griffin702/service/jwt-iris"
	"github.com/griffin702/service/tools"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

type BaseController struct {
	Ctx          iris.Context
	Session      *sessions.Session
	Svc          service.Service
	JsonPlus     model.JsonPlus
	Pager        *model.Pager
	GetClientIP  model.GetClientIP
	GetOption    model.GetOption
	Links        *model.Links
	DisableRight bool
	Hm           service.HelperMap
	Valid        model.Validator
	Tool         *tools.Tool
	Config       *config.Config
	UserID       int64
}

func (c *BaseController) BeginRequest(ctx iris.Context) {
	tokenStr := c.Session.GetString("token")
	user := c.GetUserByToken(tokenStr)
	if user.ID > 0 {
		c.UserID = user.ID
		return
	}
	tokenStr = c.Ctx.GetCookie("token")
	user = c.GetUserByToken(tokenStr)
	c.UserID = user.ID
}

func (c *BaseController) EndRequest(ctx iris.Context) {}

func (c *BaseController) IsLogin() bool {
	return c.UserID > 0
}

func (c *BaseController) GetUserByToken(tokenStr string) (user *model.UserSession) {
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
