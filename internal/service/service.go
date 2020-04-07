package service

import (
	"ginana-blog/internal/service/i_user"
	"ginana-blog/library/ecode"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type Service struct {
	db   *gorm.DB
	ef   *casbin.SyncedEnforcer
	User i_user.IUser
}

func New(
	db *gorm.DB,
	ef *casbin.SyncedEnforcer,
	u i_user.IUser,
) (s *Service, err error) {
	if err = u.SetEnforcer(ef); err != nil {
		return
	}
	s = &Service{
		db:   db,
		ef:   ef,
		User: u,
	}
	return
}

func (s *Service) Close() {
	_ = s.db.Close()
}

func (s *Service) ShowError(ctx *gin.Context, err error) {
	ec := ecode.Cause(err)
	ctx.HTML(http.StatusInternalServerError, "error/error.html", gin.H{
		"code":  ec.Code(),
		"error": ec.Message(),
	})
}
