//go:build wireinject
// +build wireinject

package main

import (
	"login-api/domains/users/handlers"
	"login-api/domains/users/repositories"
	"login-api/domains/users/usecases"
	"login-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var userSet = wire.NewSet(
	repositories.NewUserRepository,
	wire.Bind(new(repositories.UserRepositoryInterface), new(*repositories.UserRepository)),
	usecases.NewUserUseCase,
	wire.Bind(new(usecases.UserUseCaseInterface), new(*usecases.UserUseCase)),
	handlers.NewUserHttp,
	wire.Bind(new(handlers.UserHttpInterface), new(*handlers.UserHttp)),
)

func InitializeServer() (*gin.Engine, error) {
	wire.Build(
		userSet,
		routes.SetupRoutes,
	)

	return nil, nil
}
