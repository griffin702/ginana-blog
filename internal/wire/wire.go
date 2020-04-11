// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package wire

import (
	"ginana-blog/internal/config"
	"ginana-blog/internal/db"
	"ginana-blog/internal/server"
	"ginana-blog/internal/server/controller/admin"
	"ginana-blog/internal/server/controller/api"
	"ginana-blog/internal/server/controller/front"
	"ginana-blog/internal/server/router"
	"ginana-blog/internal/service"
	"ginana-blog/internal/service/i_user"
	"github.com/google/wire"
)

var initProvider = wire.NewSet(config.NewConfig, db.NewDB, db.NewCasbin)
var iProvider = wire.NewSet(i_user.New)
var cProvider = wire.NewSet(front.New, admin.New, api.New)
var httpProvider = wire.NewSet(router.InitRouter, server.NewHttpServer)

func InitApp() (*App, func(), error) {
	panic(wire.Build(
		initProvider,
		iProvider,
		service.New,
		cProvider,
		httpProvider,
		NewApp,
	))
}
