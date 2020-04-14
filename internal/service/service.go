package service

import (
	"context"
	"fmt"
	"ginana-blog/internal/config"
	"ginana-blog/library/cache/memcache"
	"ginana-blog/library/database"
	"ginana-blog/library/tools"
	"github.com/casbin/casbin/v2"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type Service interface {
	Close()
	GetError(i int, str ...string) (int, error)
	SetEnforcer(ef *casbin.SyncedEnforcer) (err error)
	GetAllRoles(ctx context.Context) (roles []database.CasbinRole, err error)
	GetAllUsers(ctx context.Context) (roles []database.CasbinUser, err error)
	GetSiteOptions() (res map[string]string, err error)
}

func New(cfg *config.Config, db *gorm.DB, mc memcache.Memcache, eh *map[int]string) (s Service, err error) {
	s = &service{
		cfg:  cfg,
		db:   db,
		mc:   mc,
		eh:   eh,
		tool: tools.New(),
	}
	_, err = s.GetSiteOptions()
	return
}

type service struct {
	cfg  *config.Config
	db   *gorm.DB
	ef   *casbin.SyncedEnforcer
	mc   memcache.Memcache
	eh   *map[int]string
	tool *tools.Tool
}

func (s *service) Close() {
	_ = s.db.Close()
}

func (s *service) GetError(i int, args ...string) (int, error) {
	if len(args) > 1 {
		panic("too many arguments")
	}
	errHelper := *s.eh
	msg := errHelper[i]
	if len(args) == 1 {
		msg = args[0]
	}
	return i, errors.New(msg)
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
