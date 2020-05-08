package admin

import (
	"ginana-blog/internal/model"
	"github.com/griffin702/service/upload"
	"strings"
)

func (c *CAdmin) PostUpload() {
	f, h, err := c.Ctx.FormFile("iNanaUploadImage")
	if err != nil {
		c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(500, err)))
		return
	}
	defer f.Close()
	albumId := c.Ctx.URLParamInt64Default("albumId", 0)
	fi, err := upload.NewFileInfo(f, h.Filename, &upload.Config{
		StaticDir:  c.Config.StaticDir,
		AlbumID:    albumId,
		LastSource: c.Ctx.URLParamDefault("last_src", ""),
		UploadType: c.Ctx.URLParamIntDefault("type", 0),
		W:          c.Ctx.URLParamIntDefault("w", 0),
		H:          c.Ctx.URLParamIntDefault("h", 0),
	})
	if err != nil {
		c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(500, err)))
		return
	}
	fi.DialogID = c.Ctx.URLParam("guid")
	path := fi.JoinInfo()
	if err := fi.SaveImage(path); err != nil {
		c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(500, err)))
		return
	}
	if albumId > 0 {
		photo := new(model.Photo)
		photo.AlbumID = albumId
		photo.Desc = fi.Name
		photo.Url = strings.TrimLeft(path, ".")
		if err = c.Svc.CreatePhoto(photo); err != nil {
			c.Ctx.JSON(c.JsonPlus(nil, err))
			return
		}
	}
	c.Ctx.JSON(c.JsonPlus(fi, c.Hm.GetMessage(0, "上传成功")))
	return
}

func (c *CAdmin) PostUploadFile() {
	f, h, err := c.Ctx.FormFile("iNanaUploadMedia")
	if err != nil {
		c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(500, err)))
		return
	}
	defer f.Close()
	fi, err := upload.NewFileInfo(f, h.Filename, &upload.Config{
		StaticDir:  c.Config.StaticDir,
		UploadType: c.Ctx.URLParamIntDefault("type", 0),
	})
	if err != nil {
		c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(500, err)))
		return
	}
	mediaPath := fi.JoinInfo()
	_, err = c.Ctx.UploadFormFiles(mediaPath)
	if err != nil {
		c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(500, err)))
		return
	}
	fi.URL = strings.TrimLeft(mediaPath, ".")
	if fi.Config.UploadType == 4 {
		if _, err := fi.GetFrame(mediaPath); err != nil {
			c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(500, err)))
			return
		}
	}
	fi.ScreenShotURL = strings.TrimLeft(fi.ScreenShotURL, ".")
	c.Ctx.JSON(c.JsonPlus(fi, c.Hm.GetMessage(0, "上传成功")))
	return
}
