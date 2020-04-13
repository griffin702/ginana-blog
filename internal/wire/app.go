package wire

import (
	"context"
	"ginana-blog/internal/service"
	"ginana-blog/library/log"
	"net/http"
	"time"
)

type App struct {
	svc  service.Service
	http *http.Server
}

func NewApp(s service.Service, h *http.Server) (app *App, closeFunc func(), err error) {
	app = &App{
		svc:  s,
		http: h,
	}
	closeFunc = func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		if err := h.Shutdown(ctx); err != nil {
			log.Errorf("httpSrv.Shutdown error(%v)", err)
		}
		s.Close()
		cancel()
	}
	return
}
