package main

import (
	"flag"
	_ "ginana-blog/docs"
	"ginana-blog/internal/config"
	"ginana-blog/internal/wire"
	"github.com/griffin702/ginana/library/conf/paladin"
	"github.com/griffin702/ginana/library/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title GiNana
// @version 1.0.0
// @description 基于GiNana的个人网站项目，默认端口：8000
// @host 127.0.0.1:8000
// @BasePath /api
// @license.name MIT License
// @license.url
func main() {
	flag.Parse()
	closeLog := log.Init()
	log.Info("GiNana App Start")
	cfg, err := config.GetBaseConfig()
	if err != nil {
		panic(err)
	}
	if err := paladin.Init(cfg.ConfigIsLocal, cfg.ConfigPath); err != nil {
		panic(err)
	}
	app, closeFunc, err := wire.InitApp()
	if err != nil {
		panic(err)
	}
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)
		for {
			s := <-ch
			log.Infof("get a signal %s", s.String())
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Info("GiNana App Exit")
				time.Sleep(time.Second)
				closeFunc()
				closeLog()
				return
			case syscall.SIGHUP:
			default:
				return
			}
		}
	}()
	err = app.Server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Errorf(err.Error())
	}
}
