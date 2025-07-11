package transaction

import (
	"context"
	"events"
	"payment-service/domains/transaction/entities"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	// Create saves a new payment record to the database within a given transaction.
	Create(ctx context.Context, tx *gorm.DB, payment *entities.Payment) error
}

// TransactionUsecase defines the business logic operations for transactions.
type TransactionUsecase interface {
	ProcessPaymentRequest(ctx context.Context, event *events.PaymentRequestedEvent) error
}
