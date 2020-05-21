package service

import (
	"context"
	"fmt"
	"ginana-blog/internal/config"
	"ginana-blog/internal/model"
	"github.com/casbin/casbin/v2"
	"github.com/griffin702/ginana/library/cache/memcache"
	"github.com/griffin702/ginana/library/database"
	"github.com/griffin702/service/tools"
	"github.com/jinzhu/gorm"
)

type Service interface {
	Close()
	SetEnforcer(ef *casbin.SyncedEnforcer) (err error)

	// 公共
	GetCaptcha() (res *model.Captcha, err error)
	PostLogin(req *model.UserLoginReq) (user *model.User, err error)
	PostRegister(req *model.UserRegisterReq) (user *model.User, err error)
	CheckPermission(userId int64, router, method string) (idAuth bool)
	GetSiteOptions() (res *model.Option, err error)
	UpdateSiteOptions(req *model.Option) (err error)
	UpdateAccount(req *model.UpdateUserReq) (user *model.User, err error)

	// 角色
	GetEFRoles(ctx context.Context) (roles []*database.EFRolePolicy, err error)
	GetRole(id int64) (role *model.Role, err error)
	GetRoles(p *model.Pager, prs ...model.RoleQueryParam) (res *model.Roles, err error)
	GetAllRoles() (res *model.Roles, err error)
	GetRoleByName(name string) (role *model.Role, err error)
	CreateRole(req *model.CreateRoleReq) (role *model.Role, err error)
	UpdateRole(req *model.UpdateRoleReq) (role *model.Role, err error)
	DeleteRole(id int64) (err error)

	// 用户
	GetEFUsers(ctx context.Context) (users []*database.EFUseRole, err error)
	GetUsers(p *model.Pager, prs ...model.UserQueryParam) (res *model.Users, err error)
	GetUser(id int64) (user *model.User, err error)
	CreateUser(req *model.CreateUserReq) (user *model.User, err error)
	UpdateUser(req *model.UpdateUserReq) (user *model.User, err error)
	DeleteUser(id int64) (err error)
	GetUserByUsername(username string) (user *model.User, err error)
	CountUsers() (count int64)

	// 规则
	GetPolicy(id int64) (policy *model.Policy, err error)
	GetPolices(p *model.Pager, prs ...model.PolicyQueryParam) (res *model.Polices, err error)
	GetAllPolices() (res *model.Polices, err error)
	CreatePolicy(req *model.CreatePolicyReq) (policy *model.Policy, err error)
	UpdatePolicy(req *model.UpdatePolicyReq) (policy *model.Policy, err error)
	DeletePolicy(id int64) (err error)

	// 文章
	GetLatestArticles(limit int) (articles []*model.Article, err error)
	GetHotArticles(limit int) (articles []*model.Article, err error)
	GetArticle(id int64) (article *model.Article, err error)
	GetArticleByUrlName(urlName string) (article *model.Article, err error)
	GetArticles(p *model.Pager, prs ...model.ArticleQueryParam) (res *model.Articles, err error)
	CreateArticle(req *model.ArticleReq) (article *model.Article, err error)
	UpdateArticle(req *model.ArticleReq) (article *model.Article, err error)
	DeleteArticle(id int64) (err error)
	BatchArticle(req *model.ArticleListReq) (err error)
	PushBaiDu(url string) (string, error)
	CountArticles() (count int64)

	// 标签
	GetTags(p *model.Pager, prs ...model.TagQueryParam) (res *model.Tags, err error)
	GetTagByName(name string) (tag *model.Tag, err error)
	BatchTag(req *model.TagListReq) (err error)
	CountTags() (count int64)
	GetTagsLimit6() (tags []*model.Tag, err error)

	// 心情
	GetMoods(p *model.Pager) (res *model.Moods, err error)
	CreateMood(req *model.MoodReq) (err error)

	// 友链
	GetLinks() (links []*model.Link, err error)
	GetLink(id int64) (link *model.Link, err error)
	CreateLink(req *model.CreateLinkReq) (link *model.Link, err error)
	UpdateLink(req *model.UpdateLinkReq) (link *model.Link, err error)
	DeleteLink(id int64) (err error)

	// 评论
	GetLatestComments(limit int) (comments []*model.Comment, err error)
	GetComments(p *model.Pager, prs ...model.CommentQueryParam) (res *model.Comments, err error)
	GetComment(id int64) (comment *model.Comment, err error)
	CreateComment(req *model.CreateCommentReq) (err error)
	UpdateComment(req *model.UpdateCommentReq) (err error)
	DeleteComment(id int64) (err error)

	// 相册
	GetAlbums(p *model.Pager, prs ...model.AlbumQueryParam) (res *model.Albums, err error)
	GetAlbum(id int64) (album *model.Album, err error)
	CreateAlbum(req *model.CreateAlbumReq) (album *model.Album, err error)
	UpdateAlbum(req *model.UpdateAlbumReq) (album *model.Album, err error)
	DeleteAlbum(id int64) (err error)
	SetAlbumStatus(id int64, hidden bool) (err error)
	SetAlbumCover(id int64, cover string) (err error)

	// 照片
	GetPhotos(p *model.Pager, albumId int64) (res *model.Photos, err error)
	CreatePhoto(req *model.CreatePhotoReq) (photo *model.Photo, err error)
	UpdatePhoto(req *model.UpdatePhotoReq) (photo *model.Photo, err error)
	DeletePhoto(id int64) (err error)
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

func (s *service) CheckPermission(userId int64, router, method string) (idAuth bool) {
	if !s.cfg.Casbin.Enable || userId == 1 {
		return true
	}
	user, err := s.GetUser(userId)
	if err != nil || !user.IsAuth {
		return
	}
	for _, role := range user.Roles {
		isAuth, err := s.ef.Enforce(role.RoleName, router, method)
		if err != nil {
			break
		}
		//fmt.Println(role.RoleName, router, method, isAuth)
		if isAuth {
			return true
		}
	}
	return
}
