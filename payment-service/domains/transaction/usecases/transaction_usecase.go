package usecases

import (
	"context"
	"errors"
	"events"
	"fmt"
	"payment-service/domains/transaction"
	"payment-service/domains/transaction/entities"
	"payment-service/infrastructures"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type TransactionUseCaseImpl struct {
	db      infrastructures.Database
	trxRepo transaction.TransactionRepository
}

func NewTransactionUsecase(db infrastructures.Database, trxRepo transaction.TransactionRepository) *TransactionUseCaseImpl {
	return &TransactionUseCaseImpl{
		db:      db,
		trxRepo: trxRepo,
	}
}

// verifyToken is a private helper method to validate the JWT from the event.
func (u *TransactionUseCaseImpl) verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the signing method is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key for validation.
		return []byte("harusnyambildarienvini"), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// ProcessPaymentRequest handles an incoming event from RabbitMQ.
func (u *TransactionUseCaseImpl) ProcessPaymentRequest(ctx context.Context, event *events.PaymentRequestedEvent) error {
	// 1. Verify the JWT from the event to ensure the message is authentic.
	token, err := u.verifyToken(event.AuthToken)
	if err != nil || !token.Valid {
		return errors.New("invalid or expired auth token in event")
	}

	// 2. Create the payment entity to be saved in this service's database.
	payment := &entities.Payment{
		ID:          event.PaymentID, // Use the ID from the event for consistency
		UserID:      event.UserID,
		Amount:      event.Amount,
		Status:      entities.PaymentPending,
		MerchantID:  event.MerchantID,
		Description: event.Description,
		CreatedAt:   event.CreatedAt, // Use the original creation time from the event
	}

	// 3. Save the pending payment record to the database.
	return u.db.GetInstance().Transaction(func(tx *gorm.DB) error {
		return u.trxRepo.Create(ctx, tx, payment) // Assuming a simple Create method in the repo
	})
}
