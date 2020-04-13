package db

import (
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
		new(model.User),
		new(model.Options),
	)
}

func initTableData(db *gorm.DB) {

}
