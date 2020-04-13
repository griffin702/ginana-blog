package admin

import (
	"ginana-blog/internal/service"
	"github.com/kataras/iris/v12"
)

type CAdmin struct {
	Ctx iris.Context
	svc service.Service
}

func New(s service.Service) *CAdmin {
	return &CAdmin{
		svc: s,
	}
}
