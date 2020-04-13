package router

import (
	"ginana-blog/internal/config"
	"ginana-blog/internal/server/controller/admin"
	"ginana-blog/internal/server/controller/api"
	"ginana-blog/internal/server/controller/front"
	"ginana-blog/internal/service"
	"ginana-blog/library/mdw"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"time"
)

func InitRouter(
	svc service.Service,
	cfg *config.Config, api *api.CApi,
	front *front.CFront, admin *admin.CAdmin,
) (e *iris.Application) {
	e = NewIris(svc, cfg)
	sessManager := sessions.New(sessions.Config{
		Cookie:  "GiNana_Session",
		Expires: 24 * time.Hour,
	})
	frontParty := mvc.New(e.Party("/"))
	frontParty.Register(sessManager.Start)
	frontParty.Router.Layout("layouts/front.html")
	frontParty.Handle(front)
	adminParty := mvc.New(e.Party("/"))
	adminParty.Router.Layout("layouts/admin.html")
	adminParty.Handle(admin)
	apiParty := e.Party("/api", mdw.Cors()).AllowMethods(iris.MethodOptions)
	{
		apiParty.Get("/", api.GetUsers)
	}
	return
}
