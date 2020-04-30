package db

import (
	"fmt"
	"ginana-blog/internal/config"
	"ginana-blog/internal/model"
	"ginana-blog/library/conf/paladin"
	"ginana-blog/library/database"
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
	initTableData(db)
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

func initTableData(db *gorm.DB) {
	admin := new(model.User)
	admin.ID = 1
	if err := db.Find(admin).Error; err != nil {
		admin.Username = "admin"
		admin.Password = "$2a$10$qhcgRHCZOsn3V8854Vw3eeJHPra.CSX4MACEIS4VqY10AazjxJxqO"
		admin.Nickname = "admin"
		admin.IsAuth = true
		db.Create(admin)
	}
	article := new(model.Article)
	if err := db.Find(article, "id = 1").Error; err != nil {
		for i := 0; i < 20; i++ {
			article = new(model.Article)
			article.UserID = 1
			article.Title = fmt.Sprintf("标题-%d", i)
			article.Content = fmt.Sprintf("内容-%d", i)
			article.Status = 1
			db.Create(article)
		}
	}
}
