package model

// Captcha
type Captcha struct {
	Image string `json:"image"` // 验证码图片
	Code  string `json:"code"`  // Code
}
