package api

import (
	"ginana-blog/internal/model"
	"ginana-blog/internal/service"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

type CApi struct {
	Ctx     iris.Context
	Session *sessions.Session
	Svc     service.Service
	Pager   *model.Pager
}

// GetLoginCaptcha godoc
// @Description 获取验证码
// @Tags Login
// @Accept  json
// @Produce  json
// @Param page query int true "页码"
// @Success 200 {object} model.Captcha
// @Failure 500 {object} model.JSON
// @Router /login/captcha [get]
func (c *CApi) GetLoginCaptcha() {
	captcha, err := c.Svc.GetCaptcha()
	if err != nil {
		c.Ctx.JSON(model.PlusJson(nil, err))
		return
	}
	c.Session.Set("Captcha", captcha.Code)
	c.Ctx.ContentType("image/png")
	c.Ctx.Write(captcha.Image)
}

func (c *CApi) GetTest() {
	code := c.Session.Get("Captcha")
	c.Ctx.JSON(model.PlusJson(code, nil))
}
