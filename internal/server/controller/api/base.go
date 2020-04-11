package api

import (
	"ginana-blog/internal/service"
	"github.com/kataras/iris/v12"
)

type CApi struct {
	Ctx iris.Context
	svc *service.Service
}

func New(s *service.Service) *CApi {
	return &CApi{
		svc: s,
	}
}
