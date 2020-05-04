package api

import (
	"ginana-blog/internal/model"
	"ginana-blog/internal/service"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

type CApi struct {
	Ctx     iris.Context
	Session *sessions.Session
	Svc     service.Service
	Pager   *model.Pager
	Hm      service.HelperMap
	Valid   model.Validator
}
