package admin

import (
	"bytes"
	"errors"
	"fmt"
	"ginana-blog/internal/model"
	"github.com/nfnt/resize"
	"github.com/ulricqin/goutils/filetool"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type Sizer interface {
	Size() int64
}

const (
	localFileDir   = "upload"
	minFileSize    = 1        // bytes
	maxFileSize    = 10000000 // bytes
	maxWidthHeight = 1280
	imageTypes     = "(jpg|gif|p?jpeg|(x-)?png)"
)

var (
	uploadTypeMap   = map[int]string{1: "bigpic", 2: "smallpic", 3: "bigsmallpic", 4: "media/mp4", 5: "media/mp3"}
	acceptFileTypes = regexp.MustCompile(imageTypes)
)

type FileInfo struct {
	Url   string `json:"url,omitempty"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Size  int64  `json:"size"`
	Error string `json:"error,omitempty"`
}

type UploadResp struct {
	URL           string `json:"url"`
	ScreenShotURL string `json:"screen_shot_url"`
	DialogID      string `json:"dialog_id"`
}

func (fi *FileInfo) ValidateType() (valid bool) {
	if acceptFileTypes.MatchString(fi.Type) {
		return true
	}
	fi.Error = "Type of file not allowed"
	return false
}

func (fi *FileInfo) ValidateSize() (valid bool) {
	if fi.Size < minFileSize {
		fi.Error = "File is too small"
	} else if fi.Size > maxFileSize {
		fi.Error = "File is too big"
	} else {
		return true
	}
	return false
}

func (c *CAdmin) PostUpload() {
	data := new(UploadResp)
	data.DialogID = c.Ctx.URLParam("guid")
	uType := c.Ctx.URLParamIntDefault("type", 1)
	f, h, err := c.Ctx.FormFile("iNanaUploadImage")
	if err != nil {
		c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(500, err)))
		return
	}
	if f != nil {
		defer f.Close()
	}
	rWidth, rHeight, fm, err := retRealWHEXT(f)
	f, h, err = c.Ctx.FormFile("iNanaUploadImage")
	if err != nil {
		c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(500, err)))
		return
	}
	fi := &FileInfo{
		Name: h.Filename,
		Type: fm,
	}
	if !fi.ValidateType() {
		c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(1010, "invalid file type")))
		return
	}
	if sizeInterface, ok := f.(Sizer); ok {
		fi.Size = sizeInterface.Size()
	}
	if !fi.ValidateSize() {
		c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(1011, fi.Error)))
		return
	}
	timeNow := time.Now().UnixNano()
	_, dir := getDefaultStaticDir(c.Config.StaticDir)
	fileSaveName := fmt.Sprintf("%s/%s/%s", localFileDir, uploadTypeMap[uType], time.Now().Format("20060102"))
	imgPath := fmt.Sprintf("%s/%s/%d.%s", dir, fileSaveName, timeNow, fm)
	_ = filetool.InsureDir(fmt.Sprintf("%s/%s", dir, fileSaveName))
	if uType == 1 { //上传类型1：文章上传的图片
		err = createSmallPicScale(f, imgPath, 0, 0, 88)
		if err != nil {
			c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(500, err)))
			return
		}
		sw, sh := retMaxWH(rWidth, rHeight, 720)
		small := changeToSmall(imgPath)
		f, _, _ = c.Ctx.FormFile("iNanaUploadImage")
		err = createSmallPicScale(f, small, sw, sh, 88)
		if err != nil {
			c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(500, err)))
			return
		}
		data.URL = strings.TrimLeft(small, ".")
		c.Ctx.JSON(c.JsonPlus(data, c.Hm.GetMessage(0, "上传成功")))
		return
	}
	if uType == 2 { //上传类型2：头像、封面等上传，只保存小图
		w := c.Ctx.URLParamIntDefault("w", 0)
		h := c.Ctx.URLParamIntDefault("h", 0)
		err = createSmallPicScale(f, imgPath, w, h, 88)
		if err != nil {
			c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(500, err)))
			return
		}
		//保存成功，则删除旧资源
		lastSrc := c.Ctx.URLParam("last_src")
		if lastSrc != "" && !c.IsDefaultSrc(lastSrc) {
			_ = os.Remove(".." + lastSrc)
		}
		data.URL = strings.TrimLeft(imgPath, ".")
		c.Ctx.JSON(c.JsonPlus(data, c.Hm.GetMessage(0, "上传成功")))
		return
	}
	if uType == 3 { //上传类型3：照片上传，同时保存大图小图
		err = createSmallPicScale(f, imgPath, 0, 0, 88)
		if err != nil {
			c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(500, err)))
			return
		}
		albumId := c.Ctx.URLParamInt64Default("albumId", 0)
		f, _, _ = c.Ctx.FormFile("iNanaUploadImage")
		if albumId > 0 {
			err = createSmallPicClip(f, changeToSmall(imgPath), rWidth, rHeight, 88)
			if err != nil {
				c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(500, err)))
				return
			}
			photo := new(model.Photo)
			photo.AlbumID = albumId
			photo.Desc = fi.Name
			photo.Url = strings.TrimLeft(imgPath, ".")
			if err = c.Svc.CreatePhoto(photo); err != nil {
				c.Ctx.JSON(c.JsonPlus(nil, err))
				return
			}
		}
		sw, sh := retMaxWH(rWidth, rHeight, 200)
		err = createSmallPicScale(f, changeToSmall(imgPath), sw, sh, 88)
		if err != nil {
			c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(500, err)))
			return
		}
		//保存成功，则删除旧资源
		lastSrc := c.Ctx.URLParam("last_src")
		if lastSrc != "" && !c.IsDefaultSrc(lastSrc) {
			_ = os.Remove(".." + lastSrc)
			_ = os.Remove(".." + changeToSmall(lastSrc))
		}
		data.URL = strings.TrimLeft(imgPath, ".")
		c.Ctx.JSON(c.JsonPlus(data, c.Hm.GetMessage(0, "上传成功")))
		return
	}
	c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(1010)))
	return
}

func (c *CAdmin) PostUploadFile() {
	data := new(UploadResp)
	uType := c.Ctx.URLParamIntDefault("type", 0)
	_, h, err := c.Ctx.FormFile("iNanaUploadMedia")
	if err != nil {
		c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(500, err)))
		return
	}
	var ext string
	list := strings.Split(h.Filename, ".")
	if len(list) > 1 {
		ext = list[1]
	}
	if uType == 4 || uType == 5 {
		if ext != "mp3" && ext != "mp4" {
			c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(1012)))
			return
		}
		timeNow := time.Now().UnixNano()
		_, dir := getDefaultStaticDir(c.Config.StaticDir)
		fileSaveName := fmt.Sprintf("%s/%s/%s", localFileDir, uploadTypeMap[uType], time.Now().Format("20060102"))
		mediaPath := fmt.Sprintf("%s/%s/%d.%s", dir, fileSaveName, timeNow, ext)
		screenShotURL := fmt.Sprintf("%s/%s/%d.jpg", dir, fileSaveName, timeNow)
		_ = filetool.InsureDir(fmt.Sprintf("%s/%s", dir, fileSaveName))
		_, err = c.Ctx.UploadFormFiles(mediaPath)
		if err != nil {
			c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(500, err)))
			return
		}
		data.URL = strings.TrimLeft(mediaPath, ".")
		if uType == 4 {
			stdout, err := getFrame(mediaPath, screenShotURL)
			if err != nil {
				c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(500, err)))
				return
			}
			log.Println(stdout)
			data.ScreenShotURL = strings.TrimLeft(screenShotURL, ".")
		}
		c.Ctx.JSON(c.JsonPlus(data, c.Hm.GetMessage(0, "上传成功")))
		return
	}
	c.Ctx.JSON(c.JsonPlus(nil, c.Hm.GetMessage(1010)))
	return
}

/*
* 图片裁剪
* 入参:1、file，2、输出路径，3、原图width，4、原图height，5、精度
* 规则:照片生成缩略图尺寸为190*135，先进行缩小，再进行平均裁剪
*
* 返回:error
 */
func createSmallPicClip(in io.Reader, fileSmall string, w, h, quality int) error {
	x0 := 0
	x1 := 190
	y0 := 0
	y1 := 135
	sh := h * 190 / w
	sw := w * 135 / h
	origin, fm, err := image.Decode(in)
	if err != nil {
		return err
	}
	if sh > 135 {
		origin = resize.Resize(uint(190), uint(sh), origin, resize.Lanczos3)
		y0 = (sh - 135) / 4
		y1 = y0 + 135
	} else {
		origin = resize.Resize(uint(sw), uint(135), origin, resize.Lanczos3)
		x0 = (sw - 190) / 2
		x1 = x0 + 190
	}
	out, err := os.Create(fileSmall)
	if err != nil {
		return err
	}
	switch fm {
	case "jpeg":
		img := origin.(*image.YCbCr)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.YCbCr)
		return jpeg.Encode(out, subImg, &jpeg.Options{Quality: quality})
	case "png":
		switch origin.(type) {
		case *image.NRGBA:
			img := origin.(*image.NRGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.NRGBA)
			return png.Encode(out, subImg)
		case *image.RGBA:
			img := origin.(*image.RGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
			return png.Encode(out, subImg)
		}
	case "gif":
		img := origin.(*image.Paletted)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.Paletted)
		return gif.Encode(out, subImg, &gif.Options{})
	default:
		return errors.New("ERROR FORMAT")
	}
	return nil
}

/*
* 缩略图生成
* 入参:1、file，2、输出路径，3、输出width，4、输出height，5、精度
* 规则: width,height是想要生成的尺寸
*
* 返回:error
 */
func createSmallPicScale(in io.Reader, fileSmall string, width, height, quality int) error {
	origin, fm, err := image.Decode(in)
	if err != nil {
		return err
	}
	if width == 0 || height == 0 {
		width = origin.Bounds().Max.X
		height = origin.Bounds().Max.Y
		//限制保存原图的长宽最大允许像素
		maxNum := maxWidthHeight
		if width < height && height > maxNum {
			width = width * maxNum / height
			height = maxNum
		} else if width >= height && width > maxNum {
			height = height * maxNum / width
			width = maxNum
		}
	}
	if quality == 0 {
		quality = 100
	}
	canvas := resize.Resize(uint(width), uint(height), origin, resize.Lanczos3)
	out, err := os.Create(fileSmall)
	if err != nil {
		return err
	}
	switch fm {
	case "jpeg":
		return jpeg.Encode(out, canvas, &jpeg.Options{Quality: quality})
	case "png":
		return png.Encode(out, canvas)
	case "gif":
		return gif.Encode(out, canvas, &gif.Options{})
	default:
		return errors.New("ERROR FORMAT")
	}
}

func changeToSmall(src string) string {
	arr1 := strings.Split(src, "/")
	filename := arr1[len(arr1)-1]
	arr2 := strings.Split(filename, ".")
	ext := "." + arr2[len(arr2)-1]
	small := strings.Replace(src, ext, "_small"+ext, 1)
	return small
}

func retMaxWH(w, h, max int) (int, int) {
	var sw, sh int
	if w < h && h > max {
		sh = max
		sw = w * max / h
	} else if w >= h && w > max {
		sw = max
		sh = h * max / w
	} else {
		sw = w
		sh = h
	}
	return sw, sh
}

func retRealWHEXT(in io.Reader) (int, int, string, error) {
	origin, fm, err := image.Decode(in)
	if err != nil {
		return 0, 0, "", err
	}
	w := origin.Bounds().Max.X
	h := origin.Bounds().Max.Y
	return w, h, fm, err
}

func getFrame(filename string, mediajpgPath string) (string, error) {
	cmd := exec.Command("ffmpeg", "-i", filename, "-y", "-f", "image2", "-t", "0.001", mediajpgPath)
	buf := new(bytes.Buffer)
	cmd.Stdout = buf
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func getDefaultStaticDir(conf string) (path, dir string) {
	staticDirList := strings.Split(conf, " ")
	if len(staticDirList) > 0 {
		def := strings.Split(staticDirList[0], ":")
		if len(def) == 2 {
			return def[0], def[1]
		}
	}
	return
}
