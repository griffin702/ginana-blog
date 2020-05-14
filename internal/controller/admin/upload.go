package admin

import (
	"ginana-blog/internal/model"
	"github.com/griffin702/service/upload"
	"github.com/kataras/iris/v12"
	"mime/multipart"
	"path/filepath"
	"strings"
)

func (c *CAdmin) PostUpload() {
	f, h, err := c.Ctx.FormFile("editormd-image-file")
	if err != nil {
		fi := new(upload.FileInfo)
		fi.Message = err.Error()
		c.Ctx.JSON(fi)
		return
	}
	defer f.Close()
	albumId := c.Ctx.URLParamInt64Default("albumId", 0)
	fi, err := upload.NewFileInfo(f, h.Filename, &upload.Config{
		StaticDir:  c.Config.StaticDir,
		AlbumID:    albumId,
		LastSource: c.Ctx.URLParamDefault("last_src", ""),
		UploadType: c.Ctx.URLParamIntDefault("type", 1),
		SmallMaxWH: c.Ctx.URLParamIntDefault("small", 0),
		W:          c.Ctx.URLParamIntDefault("w", 0),
		H:          c.Ctx.URLParamIntDefault("h", 0),
	})
	if err != nil {
		c.Ctx.JSON(fi)
		return
	}
	fi.DialogID = c.Ctx.URLParam("guid")
	path := fi.JoinInfo()
	if err := fi.SaveImage(path); err != nil {
		c.Ctx.JSON(fi)
		return
	}
	if albumId > 0 {
		photo := new(model.CreatePhotoReq)
		photo.AlbumID = albumId
		photo.Desc = fi.Name
		photo.Url = strings.TrimLeft(path, ".")
		if _, err = c.Svc.CreatePhoto(photo); err != nil {
			c.Ctx.JSON(fi)
			return
		}
	}
	c.Ctx.JSON(fi)
	return
}

func (c *CAdmin) PostUploadMedia() {
	f, h, err := c.Ctx.FormFile("filemedia")
	if err != nil {
		fi := new(upload.FileInfo)
		fi.Message = err.Error()
		c.Ctx.JSON(fi)
		return
	}
	defer f.Close()
	fi, err := upload.NewFileInfo(f, h.Filename, &upload.Config{
		StaticDir:  c.Config.StaticDir,
		UploadType: c.Ctx.URLParamIntDefault("type", 0),
	})
	if err != nil {
		c.Ctx.JSON(fi)
		return
	}
	mediaPath := fi.JoinInfo()
	_, err = c.Ctx.UploadFormFiles(filepath.Dir(mediaPath), func(ctx iris.Context, file *multipart.FileHeader) {
		file.Filename = filepath.Base(mediaPath)
	})
	if err != nil {
		c.Ctx.JSON(fi)
		return
	}
	fi.URL = strings.TrimLeft(mediaPath, ".")
	if fi.Config.UploadType == 4 {
		if _, err := fi.GetFrame(mediaPath); err != nil {
			fi.Message = err.Error()
			c.Ctx.JSON(fi)
			return
		}
	}
	fi.ScreenShotURL = strings.TrimLeft(fi.ScreenShotURL, ".")
	c.Ctx.JSON(fi)
	return
}
