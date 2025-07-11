//go:build wireinject
// +build wireinject

package main

import (
	"payment-service/domains/transaction"
	transactionrepositories "payment-service/domains/transaction/repositories"
	transactionusecases "payment-service/domains/transaction/usecases"
	"payment-service/infrastructures"
	"payment-service/infrastructures/messaging"

	"github.com/google/wire"
)

type App struct {
	Consumer *messaging.PaymentConsumer
	Usecase  transaction.TransactionUsecase
}

// NewApp is the provider function for the App struct.
func NewApp(consumer *messaging.PaymentConsumer, usecase transaction.TransactionUsecase) *App {
	return &App{Consumer: consumer, Usecase: usecase}
}

var transactionSet = wire.NewSet(
	transactionrepositories.NewTransactionRepository,
	wire.Bind(new(transaction.TransactionRepository), new(*transactionrepositories.TransactionRepositoryImpl)),
	transactionusecases.NewTransactionUsecase,
	wire.Bind(new(transaction.TransactionUsecase), new(*transactionusecases.TransactionUseCaseImpl)),
)

var databaseSet = wire.NewSet(
	infrastructures.NewPostgresDatabase,
	wire.Bind(new(infrastructures.Database), new(*infrastructures.PostgresDatabase)),
)

var messagingSet = wire.NewSet(
	messaging.NewPaymentConsumer,
)

var appSet = wire.NewSet(NewApp)

func InitializeApp() (*App, error) {
	wire.Build(
		transactionSet,
		databaseSet,
		messagingSet,
		appSet,
	)

	return nil, nil
}
