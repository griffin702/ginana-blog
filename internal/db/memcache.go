package db

import (
	"ginana-blog/internal/config"
	"ginana-blog/library/cache/memcache"
	"ginana-blog/library/conf/paladin"
)

func NewMC(cfg *config.Config) (mc memcache.Memcache, err error) {
	key := "memcache.toml"
	if err = paladin.Get(key).UnmarshalTOML(cfg); err != nil {
		return
	}
	mc = memcache.New(cfg.Memcache)
	return
}
