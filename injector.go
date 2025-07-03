//go:build wireinject
// +build wireinject

package main

import (
	accesstokens "login-api/domains/access_tokens"
	accesstokensrepositories "login-api/domains/access_tokens/repositories"
	accesstokensusecases "login-api/domains/access_tokens/usecases"
	"login-api/domains/transaction"
	transactionhandlers "login-api/domains/transaction/handlers"
	transactionrepositories "login-api/domains/transaction/repositories"
	transactionusecases "login-api/domains/transaction/usecases"
	"login-api/domains/users/handlers"
	"login-api/domains/users/repositories"
	"login-api/domains/users/usecases"
	"login-api/domains/wallet"
	wallethandlers "login-api/domains/wallet/handlers"
	walletrepositories "login-api/domains/wallet/repositories"
	walletusecases "login-api/domains/wallet/usecases"
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

var walletSet = wire.NewSet(
	walletrepositories.NewWalletRepository,
	wire.Bind(new(wallet.WalletRepository), new(*walletrepositories.WalletRepositoryImpl)),
	walletusecases.NewWalletUsecase,
	wire.Bind(new(wallet.WalletUsecase), new(*walletusecases.WalletUseCaseImpl)),
	wallethandlers.NewWalletHandler,
)

var transactionSet = wire.NewSet(
	transactionrepositories.NewTransactionRepository,
	wire.Bind(new(transaction.TransactionRepository), new(*transactionrepositories.TransactionRepositoryImpl)),
	transactionusecases.NewTransactionUsecase,
	wire.Bind(new(transaction.TransactionUsecase), new(*transactionusecases.TransactionUseCaseImpl)),
	transactionhandlers.NewTransactionHandler,
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
		walletSet,
		transactionSet,
		databaseSet,
		routes.SetupRoutes,
	)

	return nil, nil
}
