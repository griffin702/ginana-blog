package service

import (
	"ginana-blog/internal/model"
	"github.com/jinzhu/gorm"
)

func (s *service) GetCaptcha() (res *model.Captcha, err error) {
	code, image, err := s.tool.CaptchaGenerate(120, 40, 4, 0, false)
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
		return nil, s.hm.GetMessage(1001, err)
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

func (s *service) PostRegister(req *model.UserRegisterReq) (user *model.User, err error) {
	user, err = s.GetUserByUsername(req.Username)
	if err == gorm.ErrRecordNotFound {
		user = new(model.User)
		user.Username = req.Username
		user.Password = s.tool.BcryptHashGenerate(req.NewPassword)
		user.Email = req.Email
		user.Nickname = req.Nickname
		user.LastLoginIP = req.LoginIP
		user.IsAuth = true
		user.CountLogin++
		if err = s.db.Create(user).Error; err != nil {
			return nil, s.hm.GetMessage(1002, err)
		}
		return
	}
	if err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	return nil, s.hm.GetMessage(1001, "Username已被使用")
}
