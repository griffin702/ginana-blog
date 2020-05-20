package server

import (
	"ginana-blog/internal/config"
	"ginana-blog/internal/controller/admin"
	"ginana-blog/internal/controller/api"
	"ginana-blog/internal/controller/front"
	"ginana-blog/internal/model"
	"ginana-blog/internal/service"
	"github.com/griffin702/ginana/library/mdw"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

func InitRouter(svc service.Service, cfg *config.Config, hm service.HelperMap, valid model.ValidatorHandler) (e *iris.Application, err error) {

	e = newIris(svc, cfg)

	session := sessions.New(sessions.Config{
		Cookie:  "GiNana_Session",
		//Expires: time.Duration(cfg.SessionAndCookieExpire),
	})

	objects := []interface{}{
		svc, session.Start, hm, valid, getSiteOptions(svc, cfg),
		getClientIP, getPagination(svc), jsonPlus, getTools, getConfigs,
	}

	group := mvc.New(e.Party("/"))
	group.HandleError(errorHandler)
	group.Register(objects...)
	group.Router.Layout("layouts/front.html")
	group.Handle(new(front.CFront))

	adminParty := group.Party("/admin")
	adminParty.Router.Layout("layouts/admin.html")
	adminParty.Handle(new(admin.CAdmin))

	apiParty := group.Party("/api", mdw.CORS([]string{"*"}))
	apiParty.Router.AllowMethods(iris.MethodOptions)
	apiParty.Router.Layout("layouts/api.html")
	apiParty.Handle(new(api.CApiLogin))

	return
}
