package router

import (
	"ginana-blog/internal/config"
	"ginana-blog/internal/server/http/h_user"
	"github.com/gin-gonic/gin"
)

func InitRouter(u *h_user.HUser, cfg *config.Config) (e *gin.Engine) {
	e = NewGin(cfg)

	e.GET("/", u.GetUsers)

	return
}
