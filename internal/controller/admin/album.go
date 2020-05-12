package admin

import "ginana-blog/internal/model"

func (c *CAdmin) GetAlbumList() (err error) {
	albums, err := c.Svc.GetAlbums(c.Pager)
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", albums)
	c.setHeadMetas("相册列表")
	c.Ctx.View("admin/album/list.html")
	return
}

func (c *CAdmin) GetAlbumAdd() (err error) {
	c.setHeadMetas("相册创建")
	c.Ctx.View("admin/album/add.html")
	return
}

func (c *CAdmin) PostAlbumAdd() (err error) {
	req := new(model.CreateAlbumReq)
	if err = c.Ctx.ReadForm(req); err != nil {
		return
	}
	if err = c.Valid(req); err != nil {
		return
	}
	if _, err = c.Svc.CreateAlbum(req); err != nil {
		return
	}
	c.setHeadMetas("相册创建")
	c.ShowMsg("相册已创建")
	return
}

func (c *CAdmin) GetAlbumEditBy(id int64) (err error) {
	album, err := c.Svc.GetAlbum(id)
	if err != nil {
		return
	}
	c.Ctx.ViewData("data", album)
	c.setHeadMetas("相册编辑")
	c.Ctx.View("admin/album/edit.html")
	return
}

func (c *CAdmin) PostAlbumEditBy(id int64) (err error) {
	req := new(model.UpdateAlbumReq)
	if err = c.Ctx.ReadForm(req); err != nil {
		return
	}
	req.ID = id
	if err = c.Valid(req); err != nil {
		return
	}
	if _, err = c.Svc.UpdateAlbum(req); err != nil {
		return
	}
	c.setHeadMetas("相册编辑")
	c.ShowMsg("相册编辑成功")
	return
}

func (c *CAdmin) GetAlbumByShow(id int64) (err error) {
	if err = c.Svc.SetAlbumStatus(id, false); err != nil {
		return
	}
	c.setHeadMetas("显示相册")
	c.ShowMsg("设置成功")
	return
}

func (c *CAdmin) GetAlbumByHidden(id int64) (err error) {
	if err = c.Svc.SetAlbumStatus(id, true); err != nil {
		return
	}
	c.setHeadMetas("隐藏相册")
	c.ShowMsg("设置成功")
	return
}

func (c *CAdmin) GetAlbumByPhotoList(id int64) (err error) {
	photos, err := c.Svc.GetPhotos(c.Pager, id)
	if err != nil {
		return
	}
	photos.AlbumID = id
	c.Ctx.ViewData("data", photos)
	c.setHeadMetas("照片列表")
	c.Ctx.View("admin/album/photos.html")
	return
}
