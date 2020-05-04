package model

import "github.com/griffin702/service/captcha"

// Captcha
type Captcha struct {
	Image captcha.Image `json:"image"` // 验证码图片
	Code  string        `json:"code"`  // Code
}
