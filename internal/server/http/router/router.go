package router

import (
	"ginana-blog/internal/config"
	"ginana-blog/internal/server/http/h_user"
	"ginana-blog/library/log"
	"ginana-blog/library/mdw"
	"github.com/gin-gonic/gin"
)

func InitRouter(u *h_user.HUser, cfg *config.Config) (e *gin.Engine) {
	e = NewGin(cfg.GinMode)

	e.GET("/", u.GetUsers)

	return
}

func NewGin(mode string) (e *gin.Engine) {
	gin.SetMode(mode)
	gin.DefaultWriter = log.GetOutFile()
	e = gin.Default()
	e.LoadHTMLGlob("../internal/server/http/views/**/*")
	e.LoadHTMLFiles("../internal/server/http/views/index.html")
	// Logger, Recovery
	e.Use(mdw.Logger, mdw.Recovery)
	// Cors
	e.Use(mdw.CORS([]string{"*"}))
	// Swagger
	handle := mdw.SwaggerHandler("http://127.0.0.1:8000/swagger/doc.json")
	e.GET("/swagger/*any", handle)
	return
}
