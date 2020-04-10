package h_front

import (
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
	//k := &model.GiNana{
	//	Hello: "GiNana Server",
	//}
	test := `<!--[if lt IE 9]>
	<script src="/static/js/html5shiv.min.js"></script>
	<![endif]-->`
	ctx.Set("data", test)
	ctx.HTML(http.StatusOK, "front/index.html", ctx.Keys)
}
