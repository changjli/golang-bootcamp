//go:build wireinject
// +build wireinject

package main

import (
	accesstokens "login-api/domains/access_tokens"
	accesstokensrepositories "login-api/domains/access_tokens/repositories"
	accesstokensusecases "login-api/domains/access_tokens/usecases"
	"login-api/domains/users/handlers"
	"login-api/domains/users/repositories"
	"login-api/domains/users/usecases"
	"login-api/infrastructures"
	"login-api/middlewares"
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

var accessTokenSet = wire.NewSet(
	accesstokensrepositories.NewAccessTokenRepository,
	wire.Bind(new(accesstokens.AccessTokenRepositoryInterface), new(*accesstokensrepositories.AccessTokenRepository)),
	accesstokensusecases.NewAccessTokenUsecase,
	wire.Bind(new(accesstokens.AccessTokenUsecaseInterface), new(*accesstokensusecases.AccessTokenUsecase)),
)

var databaseSet = wire.NewSet(
	infrastructures.NewPostgresDatabase,
	wire.Bind(new(infrastructures.Database), new(*infrastructures.PostgresDatabase)),
)

func InitializeServer() (*gin.Engine, error) {
	wire.Build(
		middlewares.NewAuthMiddleware,
		userSet,
		accessTokenSet,
		databaseSet,
		routes.SetupRoutes,
	)

	return nil, nil
}
