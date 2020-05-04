package api

import (
	"ginana-blog/internal/model"
	"ginana-blog/internal/service"
	"github.com/kataras/iris/v12"
)

type CApi struct {
	Ctx   iris.Context
	Svc   service.Service
	Pager *model.Pager
}

// GetUsers godoc
// @Description 获取用户列表(分页)
// @Tags Users
// @Accept  json
// @Produce  json
// @Param page query int true "页码"
// @Param pagesize query int true "页码尺寸"
// @Success 200 {object} model.User
// @Failure 500 {object} resp.JSON
// @Router /users [get]
func (c *CApi) GetUsers() {
	data := model.GiNana{
		Hello: "Hello GiNana!",
	}
	c.Ctx.JSON(model.PlusJson(data, nil))
}
