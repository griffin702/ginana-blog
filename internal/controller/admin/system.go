package admin

import (
	"ginana-blog/internal/model"
)

func (c *CAdmin) GetSystemSetting() (err error) {
	c.setHeadMetas("系统设置")
	c.Ctx.View("admin/system/setting.html")
	return
}

func (c *CAdmin) PostSystemSetting() (err error) {
	option := new(model.Option)
	if err = c.Ctx.ReadForm(option); err != nil {
		return
	}
	if err = c.Svc.UpdateSiteOptions(option); err != nil {
		return
	}
	c.setHeadMetas("系统设置")
	c.ShowMsg("系统设置已更新")
	return
}
