package service

import (
	"ginana-blog/internal/model"
	"github.com/griffin702/ginana/library/ecode"
)

func (s *service) GetCaptcha() (res *model.Captcha, err error) {
	code, image, err := s.tool.CaptchaGenerate(120, 40, 4, 0)
	if err != nil {
		err = ecode.Errorf(s.hm.GetError(1005, err))
		return
	}
	res = &model.Captcha{
		Image: image,
		Code:  code,
	}
	return
}

func (s *service) PostLogin(req *model.UserLoginReq) (user *model.User, err error) {
	user, err = s.GetUserByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if !s.tool.BcryptHashCompare(user.Password, req.Password) {
		return nil, ecode.Errorf(s.hm.GetError(1008))
	}
	return
}
