package admin

import (
	"ginana-blog/internal/model"
	"os"
	"runtime"
)

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
