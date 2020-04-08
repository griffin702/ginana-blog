package http

import (
	"ginana-blog/internal/config"
	"ginana-blog/library/conf/paladin"
	"ginana-blog/library/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func NewHttpServer(e *gin.Engine, cfg *config.Config) (h *http.Server, err error) {
	if err = paladin.Get("http.toml").UnmarshalTOML(cfg); err != nil {
		return
	}
	h = &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      e,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout),
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout),
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout),
	}
	log.Printf("HTTP服务已启动 [ http://%s ]", cfg.Server.Addr)
	err = h.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Errorf(err.Error())
	}
	return
}
