package service

import (
	"context"
	"fmt"
	"ginana-blog/internal/config"
	"ginana-blog/internal/model"
	"ginana-blog/library/cache/memcache"
	"ginana-blog/library/database"
	"ginana-blog/library/tools"
	"github.com/casbin/casbin/v2"
	"github.com/jinzhu/gorm"
)

type Service interface {
	Close()
	SetEnforcer(ef *casbin.SyncedEnforcer) (err error)
	GetEFRoles(ctx context.Context) (roles []*database.EFRolePolicy, err error)
	GetEFUsers(ctx context.Context) (users []*database.EFUseRole, err error)
	GetSiteOptions() (res map[string]string, err error)

	GetLatestArticles(limit int) (articles []*model.Article, err error)
	GetHotArticles(limit int) (articles []*model.Article, err error)
	GetLatestComments(limit int) (comments []*model.Comment, err error)

	GetArticle(id int64) (article *model.Article, err error)
	GetArticles(p *model.Pager, prs ...model.ArticleQueryParam) (res *model.Articles, err error)
	GetTags() (res *model.Tags, err error)
	GetMoods(p *model.Pager) (res *model.Moods, err error)
	GetLinks() (links []*model.Link, err error)
	GetComments(p *model.Pager, objPK int64) (res *model.Comments, err error)
	GetAlbums(p *model.Pager) (res *model.Albums, err error)
	GetAlbum(id int64) (album *model.Album, err error)
	GetPhotos(p *model.Pager, albumId int64) (res *model.Photos, err error)
}

func New(cfg *config.Config, db *gorm.DB, mc memcache.Memcache, hm HelperMap) (s Service, err error) {
	s = &service{
		cfg:  cfg,
		db:   db,
		mc:   mc,
		hm:   hm,
		tool: tools.Tools,
	}
	_, err = s.GetSiteOptions()
	return
}

type service struct {
	cfg  *config.Config
	db   *gorm.DB
	ef   *casbin.SyncedEnforcer
	mc   memcache.Memcache
	hm   HelperMap
	tool *tools.Tool
}

func (s *service) Close() {
	_ = s.db.Close()
}

// Close close the resource.
func (s *service) SetEnforcer(ef *casbin.SyncedEnforcer) (err error) {
	if !s.cfg.Casbin.Enable {
		return
	}
	if s.tool.PtrIsNil(ef) {
		return fmt.Errorf("enforcer is nil")
	}
	s.ef = ef
	return
}
