package router

import (
	"fmt"
	"ginana-blog/internal/config"
	"ginana-blog/library/log"
	"ginana-blog/library/mdw"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strings"
)

func NewGin(cfg *config.Config) (e *gin.Engine) {
	gin.SetMode(cfg.GinMode)
	gin.DefaultWriter = log.GetOutFile()
	e = gin.New()
	e.Use(mdw.Logger(), mdw.Recovery())
	if cfg.EnableGzip {
		e.Use(gzip.Gzip(gzip.DefaultCompression))
	}
	initTemplate(e, cfg)
	initStaticDir(e, cfg)
	// Cors
	e.Use(mdw.CORS([]string{"*"}))
	// Swagger
	handle := mdw.SwaggerHandler("http://127.0.0.1:8000/swagger/doc.json")
	e.GET("/swagger/*any", handle)
	return
}

func initTemplate(e *gin.Engine, cfg *config.Config) {
	if !cfg.EnableTemplate {
		return
	}
	r := multitemplate.NewRenderer()
	// front
	frontLayouts, err := filepath.Glob(cfg.ViewsPath + "/layouts/front.html")
	if err != nil {
		panic(err)
	}
	frontContents, err := filepath.Glob(cfg.ViewsPath + "/front/*")
	if err != nil {
		panic(err)
	}
	// Generate our templates map from our layouts/ and includes/ directories
	for _, content := range frontContents {
		layoutCopy := make([]string, len(frontLayouts))
		copy(layoutCopy, frontLayouts)
		files := append(layoutCopy, content)
		_, dir := filepath.Split(filepath.Dir(content))
		name := fmt.Sprintf("%s/%s", dir, filepath.Base(content))
		r.AddFromFiles(name, files...)
	}
	// admin
	adminLayouts, err := filepath.Glob(cfg.ViewsPath + "/layouts/admin.html")
	if err != nil {
		panic(err)
	}
	adminContents, err := filepath.Glob(cfg.ViewsPath + "/admin/*")
	if err != nil {
		panic(err)
	}
	// Generate our templates map from our layouts/ and includes/ directories
	for _, content := range adminContents {
		layoutCopy := make([]string, len(adminLayouts))
		copy(layoutCopy, adminLayouts)
		files := append(layoutCopy, content)
		_, dir := filepath.Split(filepath.Dir(content))
		name := fmt.Sprintf("%s/%s", dir, filepath.Base(content))
		r.AddFromFiles(name, files...)
	}
	e.HTMLRender = r
}

func initStaticDir(e *gin.Engine, cfg *config.Config) {
	if !cfg.EnableTemplate {
		return
	}
	staticDirList := strings.Split(cfg.StaticDir, " ")
	if len(staticDirList) > 0 {
		path := strings.Split(staticDirList[0], ":")
		if len(path) == 2 {
			icon := "favicon.ico"
			e.StaticFile(icon, fmt.Sprintf("%s/%s", path[1], icon))
		}
	}
	for _, v := range staticDirList {
		path := strings.Split(v, ":")
		if len(path) == 2 {
			e.Static(path[0], path[1])
		}
	}
}
