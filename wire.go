//go:build wireinject

package main

import (
	"github.com/google/wire"
	"sx-go/internal/repository"
	"sx-go/internal/repository/dao"
	"sx-go/internal/service"
	"sx-go/internal/web"
	"sx-go/ioc"
)

func InitWebServer() *App {
	wire.Build(
		ioc.InitMongodb,
		ioc.InitGin,
		ioc.InitMiddlewares,
		dao.NewUserDao,
		repository.NewUserRepo,
		service.NewUserService,
		web.NewUserHandler,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
