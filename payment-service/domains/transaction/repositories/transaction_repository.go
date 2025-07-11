package repositories

import (
	"context"
	"payment-service/domains/transaction/entities"
	"payment-service/infrastructures"

	"gorm.io/gorm"
)

type TransactionRepositoryImpl struct {
	db infrastructures.Database
}

func NewTransactionRepository(db infrastructures.Database) *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{db: db}
}

func (r *TransactionRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, payment *entities.Payment) error {
	// The use case passes a transaction object 'tx', so we use that to perform the create.
	// This ensures the create operation is part of the larger transaction defined in the use case.
	return tx.WithContext(ctx).Create(payment).Error
}
