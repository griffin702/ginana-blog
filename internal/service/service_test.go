package service

import (
	"ginana-blog/internal/config"
	"ginana-blog/internal/model"
	"github.com/griffin702/ginana/library/cache/memcache"
	"github.com/griffin702/ginana/library/database"
	"github.com/jinzhu/gorm"
	"os"
	"testing"
)

var (
	svc      Service
	confPath = "../../configs/test.toml"
)

func TestMain(m *testing.M) {
	cfg, err := config.ParseToml(confPath)
	if err != nil {
		panic(err)
	}
	mc := memcache.New(cfg.Memcache)
	mysql, err := database.NewMySQL(cfg.MySQL)
	if err != nil {
		panic(err)
	}
	initTable(mysql)
	err = initTableData(mysql, mc)
	hm, err := NewHelperMap()
	if err != nil {
		panic(err)
	}
	if svc, err = New(cfg, mysql, mc, hm); err != nil {
		panic(err)
	}
	ef, err := database.NewCasbinConn(svc, cfg.ConfigPath, cfg.Casbin)
	if err != nil {
		panic(err)
	}
	if err = svc.SetEnforcer(ef); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func initTable(db *gorm.DB) {
	db.AutoMigrate(
		new(model.Options),
		new(model.User),
		new(model.Role),
		new(model.Policy),
		new(model.Article),
		new(model.Tag),
		new(model.Mood),
		new(model.Link),
		new(model.Comment),
		new(model.Album),
		new(model.Photo),
	)
}

func initTableData(db *gorm.DB, mc memcache.Memcache) (err error) {
	tx := db.Begin()
	if err = tx.Find(&model.Options{}, "id = 1").Error; err == gorm.ErrRecordNotFound {
		_ = mc.FlushAll()
		options := make(map[string]string)
		options["SiteName"] = "iNana用心交织的生活"
		options["SiteURL"] = "https://inana.top"
		options["SubTitle"] = "带着她和她去旅行"
		options["PageSize"] = "15"
		options["Keywords"] = "Python,MySQL,Golang,Windows,Linux"
		options["Description"] = "来一场说走就走的旅行"
		options["Theme"] = "main"
		options["WeiBo"] = "https://weibo.com/p/1005051484763434"
		options["Github"] = "https://github.com/griffin702"
		options["AlbumSize"] = "9"
		options["Nickname"] = "云丶先生|Nana"
		options["MyOldCity"] = "湖北省 黄石市"
		options["MyCity"] = "湖北省 武汉市"
		options["MyBirth"] = "1987-09-30"
		options["MyProfession"] = "游戏运维攻城师"
		options["MyLang"] = "Golang、Python、SQL、Shell"
		options["MyLike"] = "旅行、游戏、编程"
		options["MyWorkDesc"] = "1、Windows、Linux服务器运维，主要包括IIS、Apache、Nginx、Firewall、MySQL、SQLServer等常用服务。\r\n2、公司项目开发环境、测试环境、线上环境运维，前后端编译打包测试上线等保障工作\r\n3、日常备份与灾备恢复等确保数据安全，以及DBA相关职能。\r\n4、公司内部网络运维，硬件维护、内外网分离以及常用第三方软件运维，包括不限于SVN、FTP、Bug系统、企业邮箱等服务。\r\n5、解决不同业务需求相关各类运维脚本开发、运维工具开发、数据接口开发、Web开发等"
		for k, v := range options {
			option := new(model.Options)
			option.Name = k
			option.Value = v
			if err = tx.Create(option).Error; err != nil {
				tx.Rollback()
				return
			}
		}
	}
	tx.Commit()
	return nil
}
