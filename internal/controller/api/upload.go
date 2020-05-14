package api

import (
	"github.com/griffin702/service/upload"
	"strings"
)

// PostUploadAccountAvatar godoc
// @Description 上传当前账号头像接口
// @Tags Upload
// @Accept  json
// @Produce  json
// @Param last_src query string true "上一个资源的地址"
// @Param type query int true "上传类型，本接口只接收type=2"
// @Success 200 {string} string "{\"success\":1,\"message\":\"上传成功\",\"url\":\"xxxx/xxxx.jpg\"}"
// @Failure 500 {string} string "{\"success\":0,\"message\":\"error\",\"url\":\"\"}"
// @Router /api/upload/account/avatar [post]
func (c *CApiLogin) PostUploadAccountAvatar() {
	f, h, err := c.Ctx.FormFile("editormd-image-file")
	if err != nil {
		fi := new(upload.FileInfo)
		fi.Message = err.Error()
		c.Ctx.JSON(fi)
		return
	}
	defer f.Close()
	fi, err := upload.NewFileInfo(f, h.Filename, &upload.Config{
		StaticDir:  c.Config.StaticDir,
		LastSource: c.Ctx.URLParamDefault("last_src", ""),
		UploadType: c.Ctx.URLParamIntDefault("type", 2),
		W:          c.Ctx.URLParamIntDefault("w", 30),
		H:          c.Ctx.URLParamIntDefault("h", 30),
	})
	if err != nil {
		c.Ctx.JSON(fi)
		return
	}
	path := fi.JoinInfo()
	if err = fi.CreatePicScale(path, fi.Config.W, fi.Config.H, 88); err != nil {
		fi.Message = err.Error()
		c.Ctx.JSON(fi)
		return
	}
	if !fi.CheckSource() {
		fi.RemoveLastSource(fi.Config.LastSource, false)
	}
	fi.URL = strings.TrimLeft(path, ".")
	fi.Success = 1
	fi.Message = "上传成功"
	c.Ctx.JSON(fi)
	return
}
