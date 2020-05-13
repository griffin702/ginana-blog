package service

import (
	"ginana-blog/internal/model"
	"os"
)

func (s *service) GetAlbums(p *model.Pager, prs ...model.AlbumQueryParam) (res *model.Albums, err error) {
	var pr model.AlbumQueryParam
	if len(prs) > 0 {
		pr = prs[0]
	}
	if pr.Order == "" {
		pr.Order = "rank desc, id desc"
	}
	res = new(model.Albums)
	query := s.db.Model(&res.List).Count(&p.AllCount)
	if !pr.Admin {
		query = query.Where("hidden = 0")
	}
	query = query.Order(pr.Order).Preload("Photos")
	if err = query.Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize).Find(&res.List).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	res.Pager = p
	return
}

func (s *service) GetAlbum(id int64) (album *model.Album, err error) {
	album = new(model.Album)
	album.ID = id
	if err = s.db.Find(album).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	return
}

func (s *service) GetPhotos(p *model.Pager, albumId int64) (res *model.Photos, err error) {
	res = new(model.Photos)
	query := s.db.Model(&res.List).Where("album_id = ?", albumId).Count(&p.AllCount)
	query = query.Order("id desc")
	if err = query.Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize).Find(&res.List).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	res.Pager = p
	return
}

func (s *service) CreateAlbum(req *model.CreateAlbumReq) (album *model.Album, err error) {
	album = new(model.Album)
	album.Name = req.Name
	album.Cover = req.Cover
	album.Hidden = req.Hidden
	album.Rank = req.Rank
	if err = s.db.Create(album).Error; err != nil {
		return nil, s.hm.GetMessage(1002, err)
	}
	return
}

func (s *service) UpdateAlbum(req *model.UpdateAlbumReq) (album *model.Album, err error) {
	album = new(model.Album)
	if err = s.db.Find(album, "id = ?", req.ID).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	album.Name = req.Name
	album.Cover = req.Cover
	album.Hidden = req.Hidden
	album.Rank = req.Rank
	m, err := s.tool.StructToMap(album)
	if err != nil {
		return nil, s.hm.GetMessage(500, err)
	}
	if err = s.db.Model(album).Update(m).Error; err != nil {
		return nil, s.hm.GetMessage(1003, err)
	}
	return
}

func (s *service) SetAlbumStatus(id int64, hidden bool) (err error) {
	album := new(model.Album)
	if err = s.db.Model(album).Where("id = ?", id).Update("hidden", hidden).Error; err != nil {
		return s.hm.GetMessage(1003)
	}
	return
}

func (s *service) SetAlbumCover(id int64, cover string) (err error) {
	album := new(model.Album)
	if err = s.db.Model(album).Where("id = ?", id).Update("cover", cover).Error; err != nil {
		return s.hm.GetMessage(1003)
	}
	return
}

func (s *service) CreatePhoto(req *model.CreatePhotoReq) (photo *model.Photo, err error) {
	photo = new(model.Photo)
	photo.AlbumID = req.AlbumID
	photo.Desc = req.Desc
	photo.Url = req.Url
	if err = s.db.Create(photo).Error; err != nil {
		return nil, s.hm.GetMessage(1002, err)
	}
	return
}

func (s *service) UpdatePhoto(req *model.UpdatePhotoReq) (photo *model.Photo, err error) {
	photo = new(model.Photo)
	if err = s.db.Find(photo, "id = ?", req.ID).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	photo.AlbumID = req.AlbumID
	photo.Desc = req.Desc
	photo.Url = req.Url
	m, err := s.tool.StructToMap(photo)
	if err != nil {
		return nil, s.hm.GetMessage(500, err)
	}
	if err = s.db.Model(photo).Update(m).Error; err != nil {
		return nil, s.hm.GetMessage(1003, err)
	}
	return
}

func (s *service) DeleteAlbum(id int64) (err error) {
	album := new(model.Album)
	err = s.db.Model(album).Where("id = ?", id).Preload("Photos").Find(album).Error
	if err != nil {
		return s.hm.GetMessage(1001, err)
	}
	photoMap := make(map[string]string)
	tx := s.db.Begin()
	for _, photo := range album.Photos {
		photoMap[photo.Url] = photo.ChangetoSmall()
		if err = tx.Delete(photo).Error; err != nil {
			tx.Rollback()
			return s.hm.GetMessage(1004, err)
		}
	}
	if err = tx.Delete(album).Error; err != nil {
		tx.Rollback()
		return s.hm.GetMessage(1004, err)
	}
	tx.Commit()
	for url, small := range photoMap {
		s.deletePhotoSource(url, small)
	}
	return
}

func (s *service) DeletePhoto(id int64) (err error) {
	photo := new(model.Photo)
	if err = s.db.Delete(photo, "id = ?", id).Error; err != nil {
		return s.hm.GetMessage(1004, err)
	}
	s.deletePhotoSource(photo.Url, photo.ChangetoSmall())
	return
}

func (s *service) deletePhotoSource(url, small string) {
	_ = os.Remove(".." + url)
	_ = os.Remove(".." + small)
}
