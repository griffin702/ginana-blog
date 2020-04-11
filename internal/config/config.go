package config

import (
	"context"
	"ginana-blog/library/conf/paladin"
	"ginana-blog/library/database"
	"ginana-blog/library/log"
	xtime "ginana-blog/library/time"
	"github.com/BurntSushi/toml"
)

var (
	global  *Config
	defPath = "../configs/global.toml" // 默认配置文件
)

func GetBaseConfig() (cfg *Config, err error) {
	cfg, err = ParseToml(defPath)
	global = cfg
	return
}

// ParseToml 解析toml配置文件
func ParseToml(path string) (*Config, error) {
	var c Config
	_, err := toml.DecodeFile(path, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// NewConfig 初始化全局配置并自动载入
func NewConfig() (cfg *Config, err error) {
	key := "global.toml"
	if err = paladin.Get(key).UnmarshalTOML(global); err != nil {
		return
	}
	cfg = global
	go func() {
		for range paladin.WatchEvent(context.Background(), key) {
			if err := paladin.Get(key).UnmarshalTOML(global); err != nil {
				log.Errorf("config load error: %v", err)
				continue
			}
		}
	}()
	return
}

// Global 获取全局配置
func Global() *Config {
	if global == nil {
		return &Config{}
	}
	return global
}

type Config struct {
	AppName        string
	Version        string
	ConfigIsLocal  bool
	ConfigPath     string
	MySQL          *database.SQLConfig
	Casbin         *database.CasbinConfig
	Server         *ServerConfig
	CacheExpire    xtime.Duration
	IrisLogLevel   string
	EnableGzip     bool
	EnableTemplate bool
	ReloadTemplate bool
	ViewsPath      string
	StaticDir      string
}

type ServerConfig struct {
	Addr         string
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
	IdleTimeout  xtime.Duration
}
