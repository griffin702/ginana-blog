package public

import (
	"ginana-blog/internal/controller"
	"github.com/kataras/iris/v12"
)

type CPublic struct {
	controller.BaseController
}

// 重写BeginRequest 处理未登录时重定向到CFront
func (c *CPublic) BeginRequest(ctx iris.Context) {
	user := c.ParseToken()
	c.UserID = user.ID
	if user.ID > 0 && c.Ctx.Path() != "/public/logout" { // 已登陆
		ctx.Redirect("/admin")
		return
	}
}
