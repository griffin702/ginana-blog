package router

import (
	"ginana/internal/server/http/h_user"
	"ginana/library/log"
	"ginana/library/mdw"
	"github.com/gin-gonic/gin"
)

func InitRouter(u *h_user.HUser) (e *gin.Engine) {
	e = NewGin()

	e.GET("/", u.GetUsers)

	return
}

func NewGin() (e *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = log.GetOutFile()
	e = gin.Default()
	// Logger, Recovery
	e.Use(mdw.Logger, mdw.Recovery)
	// Cors
	e.Use(mdw.CORS([]string{"*"}))
	// Swagger
	handle := mdw.SwaggerHandler("http://127.0.0.1:8000/swagger/doc.json")
	e.GET("/swagger/*any", handle)
	return
}
