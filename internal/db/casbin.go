package db

import (
	"context"
	"ginana-blog/internal/config"
	"ginana-blog/internal/service/i_user"
	"ginana-blog/library/conf/paladin"
	"ginana-blog/library/database"
	"ginana-blog/library/log"
	"github.com/casbin/casbin/v2"
	"time"
)

// NewCasbin after Service, and SetEnforcer for service
func NewCasbin(user i_user.IUser, cfg *config.Config) (ef *casbin.SyncedEnforcer, err error) {
	key := "casbin.toml"
	if err = paladin.Get(key).UnmarshalTOML(cfg); err != nil {
		return
	}
	ef, err = database.NewCasbinConn(user, cfg.ConfigPath, cfg.Casbin)
	if err != nil {
		return
	}
	go WatchCasbinModel(ef, cfg.Casbin)
	go WatchCasbinConfig(ef, cfg.Casbin, key)
	return
}

func WatchCasbinModel(e *casbin.SyncedEnforcer, c *database.CasbinConfig) {
	for range paladin.WatchEvent(context.Background(), c.Model) {
		if err := e.LoadModel(); err != nil {
			log.Printf("e.LoadModel error(%v)", err)
		}
	}
}

func WatchCasbinConfig(e *casbin.SyncedEnforcer, c *database.CasbinConfig, key string) {
	for event := range paladin.WatchEvent(context.Background(), key) {
		autoLoad := c.AutoLoad
		autoLoadInternal := c.AutoLoadInternal
		s := &paladin.TOML{}
		_ = s.Set(event.Value)
		if err := s.Get("Casbin").UnmarshalTOML(c); err != nil {
			continue
		}
		if c.AutoLoad != autoLoad {
			if c.AutoLoad {
				e.StartAutoLoadPolicy(time.Duration(c.AutoLoadInternal) * time.Second)
			} else {
				e.StopAutoLoadPolicy()
			}
		}
		if c.AutoLoadInternal != autoLoadInternal {
			e.StopAutoLoadPolicy()
			e.StartAutoLoadPolicy(time.Duration(c.AutoLoadInternal) * time.Second)
		}
		if !c.Enable {
			e.StopAutoLoadPolicy()
		}
	}
}
