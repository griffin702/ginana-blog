package api

import (
	"ginana-blog/internal/model"
	"ginana-blog/internal/server/resp"
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

func (c *CApi) GetUsers(ctx iris.Context) {
	data := model.GiNana{
		Hello: "Hello GiNana!",
	}
	ctx.JSON(resp.PlusJson(data, nil))
}
