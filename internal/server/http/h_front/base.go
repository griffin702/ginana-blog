package h_front

import (
	"ginana-blog/internal/model"
	"ginana-blog/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HFront struct {
	svc *service.Service
}

func New(s *service.Service) *HFront {
	return &HFront{
		svc: s,
	}
}

func (h *HFront) Index(ctx *gin.Context) {
	k := &model.GiNana{
		Hello: "GiNana Server",
	}
	ctx.HTML(http.StatusOK, "front/index.html", k)
}
