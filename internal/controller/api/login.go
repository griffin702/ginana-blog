package api

import (
	"ginana-blog/internal/model"
	"github.com/griffin702/ginana/library/ecode"
)

// GetLoginCaptcha godoc
// @Description 获取验证码
// @Tags Login
// @Accept  json
// @Produce  json
// @Success 200 {string} string "data:image/png;base64,U3dhZ2dlciByb2Nrcw=="
// @Failure 500 {object} model.JSON
// @Router /login/captcha [get]
func (c *CApi) GetLoginCaptcha() {
	captcha, err := c.Svc.GetCaptcha()
	if err != nil {
		c.Ctx.JSON(model.PlusJson(nil, err))
		return
	}
	c.Session.Set(c.Hm.GetCacheKey(8), captcha.Code)
	c.Ctx.ContentType("image/png")
	c.Ctx.Write(captcha.Image)
}

// PostLoginCaptchaCheck godoc
// @Description 提前检查验证码
// @Tags Login
// @Accept  json
// @Produce  json
// @Param captcha body model.Captcha true "Check Captcha"
// @Success 200 {object} model.JSON{data=bool}
// @Failure 500 {object} model.JSON
// @Router /login/captcha/check [post]
func (c *CApi) PostLoginCaptchaCheck() {
	captcha := new(model.Captcha)
	if err := c.Ctx.ReadJSON(&captcha); err != nil {
		c.Ctx.JSON(model.PlusJson(nil, err))
		return
	}
	code := c.Session.Get(c.Hm.GetCacheKey(8))
	if code == nil {
		err := ecode.Errorf(c.Hm.GetError(1006))
		c.Ctx.JSON(model.PlusJson(nil, err))
		return
	}
	if captcha.Code == code {
		c.Ctx.JSON(model.PlusJson(true, nil))
		return
	}
	c.Ctx.JSON(model.PlusJson(false, nil))
	return
}

func (c *CApi) PostLogin() {
	req := new(model.UserLoginReq)
	if err := c.Ctx.ReadJSON(req); err != nil {
		c.Ctx.JSON(model.PlusJson(nil, err))
		return
	}
	if err := c.Valid(req); err != nil {
		c.Ctx.JSON(model.PlusJson(nil, err))
		return
	}
	code := c.Session.Get(c.Hm.GetCacheKey(8))
	if code == nil {
		err := ecode.Errorf(c.Hm.GetError(1006))
		c.Ctx.JSON(model.PlusJson(nil, err))
		return
	}
	if code != req.Code {
		err := ecode.Errorf(c.Hm.GetError(1007))
		c.Ctx.JSON(model.PlusJson(nil, err))
		return
	}
}
