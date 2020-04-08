package h_admin

import (
	"ginana-blog/internal/model"
	"ginana-blog/internal/service"
	"github.com/gin-gonic/gin"
)

type HAdmin struct {
	svc *service.Service
}

func New(s *service.Service) *HAdmin {
	return &HAdmin{
		svc: s,
	}
}

func (h *HAdmin) AdminIndex(ctx *gin.Context) {
	k := &model.GiNana{
		Hello: "GiNana Server",
	}
	ctx.HTML(200, "admin/index.html", k)
}
