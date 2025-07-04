//go:build wireinject
// +build wireinject

package main

import (
	accesstokens "payment-service/domains/access_tokens"
	accesstokensrepositories "payment-service/domains/access_tokens/repositories"
	accesstokensusecases "payment-service/domains/access_tokens/usecases"
	"payment-service/domains/transaction"
	transactionhandlers "payment-service/domains/transaction/handlers"
	transactionrepositories "payment-service/domains/transaction/repositories"
	transactionusecases "payment-service/domains/transaction/usecases"
	"payment-service/domains/users/handlers"
	"payment-service/domains/users/repositories"
	"payment-service/domains/users/usecases"
	walletservice "payment-service/domains/wallet_service"
	walletclients "payment-service/domains/wallet_service/clients"
	"payment-service/infrastructures"
	"payment-service/middlewares"
	"payment-service/routes"
	"payment-service/wizards"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var configSet = wire.NewSet(
	wizards.NewConfig,
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
	walletclients.NewWalletServiceClient,
	wire.Bind(new(walletservice.WalletServiceClient), new(*walletclients.WalletServiceClientImpl)),
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
		configSet,
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
