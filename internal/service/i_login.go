package service

import (
	"ginana-blog/internal/model"
)

func (s *service) GetCaptcha() (res *model.Captcha, err error) {
	code, image, err := s.tool.CaptchaGenerate(120, 40, 4, 0)
	if err != nil {
		return nil, s.hm.GetMessage(1005, err)
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
		return nil, s.hm.GetMessage(1008)
	}
	if !user.IsAuth {
		return nil, s.hm.GetMessage(1009)
	}
	user.CountLogin++
	user.LastLoginIP = req.LoginIP
	if err = s.db.Model(user).Select("last_login_ip", "count_login").
		Update(user).Error; err != nil {
		return nil, s.hm.GetMessage(1003, err)
	}
	return
}
