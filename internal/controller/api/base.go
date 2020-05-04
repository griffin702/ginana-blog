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
// @Success 200 []byte "image/png"
// @Failure 500 {object} model.JSON
// @Router /login/captcha [get]
func (c *CApi) GetLoginCaptcha() {
	captcha, err := c.Svc.GetCaptcha()
	if err != nil {
		c.Ctx.JSON(model.PlusJson(nil, err))
		return
	}
	c.Session.Set("captchaCode", captcha.Code)
	c.Ctx.ContentType("image/png")
	c.Ctx.Write(captcha.Image)
}

func (c *CApi) GetTest() {
	code := c.Session.Get("Captcha")
	c.Ctx.JSON(model.PlusJson(code, nil))
}

// PostLoginCaptchaCheck godoc
// @Description 提前检查验证码
// @Tags Login
// @Accept  json
// @Produce  json
// @Param code body model.Captcha true "Check Captcha"
// @Success 200 bool
// @Failure 500 {object} model.JSON
// @Router /login/captcha/check [post]
func (c *CApi) PostLoginCaptchaCheck() {
	captcha := new(model.Captcha)
	if err := c.Ctx.ReadJSON(&captcha); err != nil {
		c.Ctx.JSON(model.PlusJson(nil, err))
		return
	}
	code := c.Session.Get("captchaCode")
	if code == nil {
		c.Ctx.JSON(model.PlusJson(nil, "验证码不存在"))
		return
	}
	if captcha.Code == code {
		c.Ctx.JSON(model.PlusJson(true, nil))
		return
	}
	c.Ctx.JSON(model.PlusJson(false, nil))
	return
}
