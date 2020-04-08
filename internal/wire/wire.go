// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package wire

import (
	"ginana-blog/internal/config"
	"ginana-blog/internal/db"
	"ginana-blog/internal/server/http"
	"ginana-blog/internal/server/http/h_admin"
	"ginana-blog/internal/server/http/h_api"
	"ginana-blog/internal/server/http/h_front"
	"ginana-blog/internal/server/http/router"
	"ginana-blog/internal/service"
	"ginana-blog/internal/service/i_user"
	"github.com/google/wire"
)

var initProvider = wire.NewSet(config.NewConfig, db.NewDB, db.NewCasbin)
var iProvider = wire.NewSet(i_user.New)
var hProvider = wire.NewSet(h_front.New, h_admin.New, h_api.New)
var httpProvider = wire.NewSet(router.InitRouter, http.NewHttpServer)

func InitApp() (*App, func(), error) {
	panic(wire.Build(
		initProvider,
		iProvider,
		service.New,
		hProvider,
		httpProvider,
		NewApp,
	))
}
