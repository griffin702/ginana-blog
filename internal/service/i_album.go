package service

import (
	"ginana-blog/internal/model"
	"github.com/griffin702/ginana/library/ecode"
)

func (s *service) GetAlbums(p *model.Pager) (res *model.Albums, err error) {
	res = new(model.Albums)
	query := s.db.Model(&res.List).Count(&p.AllCount)
	query = query.Order("created_at desc").Preload("Photos")
	if err = query.Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize).Find(&res.List).Error; err != nil {
		err = ecode.Errorf(s.hm.GetError(1001, err))
		return nil, err
	}
	res.Pager = p.NewPager(p.UrlPath)
	return
}

func (s *service) GetAlbum(id int64) (album *model.Album, err error) {
	album = new(model.Album)
	album.ID = id
	if err = s.db.Find(album).Error; err != nil {
		err = ecode.Errorf(s.hm.GetError(1001, err))
		return nil, err
	}
	return
}

func (s *service) GetPhotos(p *model.Pager, albumId int64) (res *model.Photos, err error) {
	res = new(model.Photos)
	query := s.db.Model(&res.List).Where("album_id = ?", albumId).Count(&p.AllCount)
	query = query.Order("created_at desc")
	if err = query.Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize).Find(&res.List).Error; err != nil {
		err = ecode.Errorf(s.hm.GetError(1001, err))
		return nil, err
	}
	res.Pager = p.NewPager(p.UrlPath)
	return
}
