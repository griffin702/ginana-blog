package db

import (
	"fmt"
	"ginana-blog/internal/config"
	"ginana-blog/internal/model"
	"github.com/griffin702/ginana/library/conf/paladin"
	"github.com/griffin702/ginana/library/database"
	"github.com/jinzhu/gorm"
)

func NewDB(cfg *config.Config) (db *gorm.DB, err error) {
	key := "db.toml"
	if err = paladin.Get(key).UnmarshalTOML(cfg); err != nil {
		return
	}
	db, err = database.NewMySQL(cfg.MySQL)
	if err != nil {
		return
	}
	if cfg.MySQL.Debug {
		db = db.Debug()
	}
	initTable(db)
	err = initTableData(db)
	return
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

func initTableData(db *gorm.DB) (err error) {
	tx := db.Begin()
	admin := new(model.User)
	if err = tx.Find(admin, "id = 1").Error; err != nil {
		admin.Username = "admin"
		admin.Password = "$2a$10$qhcgRHCZOsn3V8854Vw3eeJHPra.CSX4MACEIS4VqY10AazjxJxqO"
		admin.Nickname = "admin"
		admin.IsAuth = true
		if err = tx.Create(admin).Error; err != nil {
			tx.Rollback()
			return
		}
	}
	if err = tx.Find(&model.Options{}, "id = 1").Error; err != nil {
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
	link := new(model.Link)
	if err = tx.Find(link, "id = 1").Error; err != nil {
		link.SiteName = "iNana"
		link.SiteAvatar = "/static/upload/default/user-default-60x60.png"
		link.SiteDesc = "iNana个人博客"
		link.Url = "https://www.inana.top"
		link.Rank = 100
		if err = tx.Create(link).Error; err != nil {
			tx.Rollback()
			return
		}
		link = new(model.Link)
		link.SiteName = "爱在发烧"
		link.SiteAvatar = "/static/upload/default/user-default-60x60.png"
		link.SiteDesc = "一个非常棒的站点，博主也很厉害"
		link.Url = "http://azfashao.com"
		link.Rank = 99
		if err = tx.Create(link).Error; err != nil {
			tx.Rollback()
			return
		}
	}
	if err = tx.Find(&model.Article{}, "id = 1").Error; err != nil {
		for i := 0; i < 20; i++ {
			article := new(model.Article)
			article.UserID = 1
			article.Title = fmt.Sprintf("标题-%d", i)
			article.Content = fmt.Sprintf("内容-%d", i)
			tag := new(model.Tag)
			tag.Name = fmt.Sprintf("标签-%d", i)
			if err = tx.Create(tag).Error; err != nil {
				tx.Rollback()
				return
			}
			article.Tags = append(article.Tags, tag)
			if err = tx.Create(article).Error; err != nil {
				tx.Rollback()
				return
			}
		}
	}
	tx.Commit()
	return nil
}
