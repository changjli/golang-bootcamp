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
	"core-service/domains/users"
	userhandlers "core-service/domains/users/handlers"
	userrepositories "core-service/domains/users/repositories"
	userusecases "core-service/domains/users/usecases"
	"core-service/domains/wallet"
	wallethandlers "core-service/domains/wallet/handlers"
	walletrepositories "core-service/domains/wallet/repositories"
	walletusecases "core-service/domains/wallet/usecases"
	"core-service/infrastructures"
	"core-service/infrastructures/messaging"
	"core-service/middlewares"
	"core-service/routes"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var userSet = wire.NewSet(
	userrepositories.NewUserRepository,
	wire.Bind(new(users.UserRepository), new(*userrepositories.UserRepositoryImpl)),
	userusecases.NewUserUseCase,
	wire.Bind(new(users.UserUseCase), new(*userusecases.UserUseCaseImpl)),
	userhandlers.NewUserHttp,
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

var messagingSet = wire.NewSet(
	messaging.NewRabbitMQPublisher,
	wire.Bind(new(messaging.PaymentPublisher), new(*messaging.RabbitMqPublisher)),
)

func InitializeServer() (*gin.Engine, error) {
	wire.Build(
		middlewares.NewAuthMiddleware,
		messagingSet,
		userSet,
		accessTokenSet,
		walletSet,
		transactionSet,
		databaseSet,
		routes.SetupRoutes,
	)

	return nil, nil
}

func InitializeMigrator() (*infrastructures.PostgresDatabase, error) {
	wire.Build(
		databaseSet,
	)
	return nil, nil
}
