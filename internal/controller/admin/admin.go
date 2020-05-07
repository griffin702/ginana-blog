package admin

import (
	"fmt"
	"ginana-blog/internal/controller"
	"ginana-blog/internal/model"
	"os"
	"runtime"
)

type CAdmin struct {
	controller.BaseController
}

func (c *CAdmin) setHeadMetas(params ...string) {
	c.Ctx.ViewData("isLogin", c.IsLogin())
	title := fmt.Sprintf("inana 后台管理 v%s", c.Config.Version)
	if len(params) > 0 {
		title = fmt.Sprintf("%s - %s", params[0], title)
	}
	c.Ctx.ViewData("title", title)
}

func (c *CAdmin) Get() (err error) {
	adminData := new(model.AdminData)
	adminData.Hostname, _ = os.Hostname()
	adminData.Gover = runtime.Version()
	adminData.OS = runtime.GOOS
	adminData.CountCpu = runtime.NumCPU()
	adminData.Arch = runtime.GOARCH
	adminData.CountArticles = c.Svc.CountArticles()
	adminData.CountUsers = c.Svc.CountUsers()
	adminData.CountTags = c.Svc.CountTags()
	c.Ctx.ViewData("data", adminData)
	c.setHeadMetas()
	c.Ctx.View("admin/index.html")
	return
}
