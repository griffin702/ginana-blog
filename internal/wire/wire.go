// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package wire

import (
	"ginana-blog/internal/config"
	"ginana-blog/internal/db"
	"ginana-blog/internal/server"
	"ginana-blog/internal/service"
	"github.com/google/wire"
)

var initProvider = wire.NewSet(config.NewConfig, db.NewDB, db.NewMC)
var svcProvider = wire.NewSet(service.NewHelperMap, service.New, db.NewCasbin)
var httpProvider = wire.NewSet(server.InitRouter, server.NewHttpServer)

func InitApp() (*App, func(), error) {
	panic(wire.Build(
		initProvider,
		svcProvider,
		httpProvider,
		NewApp,
	))
}
