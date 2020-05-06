package front

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

type CFront struct {
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

func (c *CFront) BeginRequest(ctx iris.Context) {
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

func (c *CFront) EndRequest(ctx iris.Context) {}

func (c *CFront) IsLogin() bool {
	return c.UserID > 0
}

func (c *CFront) setHeadMetas(params ...string) {
	c.Ctx.ViewData("isLogin", c.IsLogin())
	c.Ctx.ViewData("disableRight", c.DisableRight)
	titleBuf := make([]string, 0, 3)
	if len(params) == 0 && c.GetOption("sitename") != "" {
		titleBuf = append(titleBuf, c.GetOption("sitename"))
	}
	if len(params) > 0 {
		titleBuf = append(titleBuf, params[0])
	}
	titleBuf = append(titleBuf, c.GetOption("subtitle"))
	c.Ctx.ViewData("title", strings.Join(titleBuf, " - "))
	if len(params) > 1 {
		c.Ctx.ViewData("keywords", params[1])
	} else {
		c.Ctx.ViewData("keywords", c.GetOption("keywords"))
	}
	if len(params) > 2 {
		c.Ctx.ViewData("description", params[2])
	} else {
		c.Ctx.ViewData("description", c.GetOption("description"))
	}
}

func (c *CFront) getUserByToken(tokenStr string) (user *model.UserSession) {
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
