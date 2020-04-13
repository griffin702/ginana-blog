package front

import (
	"ginana-blog/internal/model"
	"ginana-blog/internal/service"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

type CFront struct {
	Ctx     iris.Context
	Session *sessions.Session
	svc     service.Service
}

func New(s service.Service) *CFront {
	return &CFront{
		svc: s,
	}
}

func (c *CFront) Get() {
	data := model.GiNana{
		Hello: "Hello GiNana!",
	}
	c.Ctx.ViewData("data", data)
	c.Ctx.View("front/index.html")
}
