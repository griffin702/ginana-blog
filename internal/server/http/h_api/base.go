package h_api

import (
	"ginana-blog/internal/model"
	"ginana-blog/internal/service"
	"github.com/gin-gonic/gin"
)

type HApi struct {
	svc *service.Service
}

func New(s *service.Service) *HApi {
	return &HApi{
		svc: s,
	}
}

// GetUsers godoc
// @Description 测试
// @Tags Public
// @Accept  json
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {string} string "ok"
// @Failure 500 {string} string "failed"
// @Router /users [get]
func (h *HApi) GetUsers(ctx *gin.Context) {
	k := &model.GiNana{
		Hello: "GiNana Server",
	}
	ctx.JSON(200, k)
}
