package server

import (
	"fmt"
	"ginana-blog/internal/config"
	"ginana-blog/library/conf/paladin"
	"ginana-blog/library/log"
	"ginana-blog/library/mdw"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"net/http"
	"strings"
	"time"
)

func NewHttpServer(irisApp *iris.Application, cfg *config.Config) (h *http.Server, err error) {
	if err = paladin.Get("http.toml").UnmarshalTOML(cfg); err != nil {
		return
	}
	if err = irisApp.Build(); err != nil {
		log.Println(err.Error())
	}
	h = &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      irisApp,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout),
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout),
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout),
	}
	log.Printf("HTTP服务已启动 [ http://%s ]", cfg.Server.Addr)
	return
}

func newIris(cfg *config.Config) (e *iris.Application) {
	e = iris.New()
	//e.Use(iris.Cache304(10 * time.Second))
	golog.Install(log.GetLogger())
	customLogger := logger.New(logger.Config{
		Status: true, IP: true, Method: true, Path: true, Query: true,
		//MessageHeaderKeys: []string{"User-Agent"},
	})
	e.OnAnyErrorCode(customLogger)
	e.Use(customLogger, recover.New())
	e.Logger().SetLevel(cfg.IrisLogLevel)
	initTemplate(e, cfg)
	initStaticDir(e, cfg)
	// swagger
	handle := mdw.SwaggerHandler("http://127.0.0.1:8000/swagger/doc.json")
	e.Get("/swagger/*any", handle)

	e.Use(func(ctx iris.Context) {
		ctx.Gzip(cfg.EnableGzip)
		ctx.Next()
	})
	return
}

func initStaticDir(e *iris.Application, cfg *config.Config) {
	if !cfg.EnableTemplate {
		return
	}
	staticDirList := strings.Split(cfg.StaticDir, " ")
	if len(staticDirList) > 0 {
		path := strings.Split(staticDirList[0], ":")
		e.Favicon(fmt.Sprintf("%s/favicon.ico", path[1]))
	}
	for _, v := range staticDirList {
		path := strings.Split(v, ":")
		if len(path) == 2 {
			e.HandleDir(path[0], path[1], iris.DirOptions{
				Gzip:     true,
				ShowList: false,
				//IndexName: "/index.html",
			})
		}
	}
	return
}

func getDefaultStaticDir(conf string) (path, dir string) {
	staticDirList := strings.Split(conf, " ")
	if len(staticDirList) > 0 {
		def := strings.Split(staticDirList[0], ":")
		if len(def) == 2 {
			return def[0], def[1]
		}
	}
	return
}
