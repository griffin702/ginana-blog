package router

import (
	"ginana-blog/internal/config"
	"ginana-blog/internal/server/http/h_admin"
	"ginana-blog/internal/server/http/h_api"
	"ginana-blog/internal/server/http/h_front"
	"github.com/gin-gonic/gin"
)

func InitRouter(front *h_front.HFront, admin *h_admin.HAdmin, api *h_api.HApi, cfg *config.Config) (e *gin.Engine) {
	e = NewGin(cfg)
	fr := e.Group("/")
	{
		fr.GET("/", front.Index)
	}
	ad := e.Group("/admin")
	{
		ad.GET("/admin", admin.AdminIndex)
	}
	ap := e.Group("/api")
	{
		ap.GET("/users", api.GetUsers)
	}
	return
}
