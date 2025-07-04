//go:build wireinject
// +build wireinject

package main

import (
	accesstokens "core-service/domains/access_tokens"
	accesstokensrepositories "core-service/domains/access_tokens/repositories"
	accesstokensusecases "core-service/domains/access_tokens/usecases"
	"core-service/domains/transaction"
	transactionhandlers "core-service/domains/transaction/handlers"
	transactionrepositories "core-service/domains/transaction/repositories"
	transactionusecases "core-service/domains/transaction/usecases"
	"core-service/domains/users/handlers"
	"core-service/domains/users/repositories"
	"core-service/domains/users/usecases"
	"core-service/domains/wallet"
	wallethandlers "core-service/domains/wallet/handlers"
	walletrepositories "core-service/domains/wallet/repositories"
	walletusecases "core-service/domains/wallet/usecases"
	"core-service/infrastructures"
	"core-service/middlewares"
	"core-service/routes"

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
