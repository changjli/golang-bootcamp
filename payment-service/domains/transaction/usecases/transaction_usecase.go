package usecases

import (
	"context"
	"fmt"
	"time"

	"payment-service/domains/transaction"
	"payment-service/domains/transaction/entities"
	"payment-service/domains/transaction/models/requests"
	"payment-service/domains/transaction/models/responses"
	walletservice "payment-service/domains/wallet_service"
	"payment-service/infrastructures"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionUseCaseImpl struct {
	db                  infrastructures.Database
	trxRepo             transaction.TransactionRepository
	walletServiceClient walletservice.WalletServiceClient
}

func NewTransactionUsecase(db infrastructures.Database, trxRepo transaction.TransactionRepository, walletServiceClient walletservice.WalletServiceClient) *TransactionUseCaseImpl {
	return &TransactionUseCaseImpl{
		db:                  db,
		trxRepo:             trxRepo,
		walletServiceClient: walletServiceClient,
	}
}

func (u *TransactionUseCaseImpl) InitiatePayment(ctx *gin.Context, userID string, req *requests.PayRequest) (*responses.PayResponse, error) {
	err := u.walletServiceClient.VerifyBalance(ctx, userID, req.Amount)
	if err != nil {
		return nil, err
	}

	paymentTrx := &entities.Transaction{
		ID:          uuid.NewString(),
		UserID:      userID,
		Type:        entities.Payment,
		Amount:      req.Amount,
		Status:      entities.Pending,
		MerchantID:  req.MerchantID,
		Description: req.Description,
	}

	err = u.db.GetInstance().Transaction(func(tx *gorm.DB) error {
		return u.trxRepo.CreateInTx(ctx, tx, paymentTrx)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}

	response := &responses.PayResponse{
		TransactionID: paymentTrx.ID,
		Message:       "Payment initiated successfully. Awaiting confirmation.",
		NewBalance:    0,
	}
	return response, nil
}

func (u *TransactionUseCaseImpl) ExpirePayments(ctx context.Context) error {
	expirationTime := time.Now().Add(-10 * time.Minute)
	expiredTrx, err := u.trxRepo.FindPendingPaymentsBefore(ctx, expirationTime)
	if err != nil {
		return fmt.Errorf("failed to find expired payments: %w", err)
	}

	if len(expiredTrx) == 0 {
		return nil
	}

	idsToExpire := make([]string, len(expiredTrx))
	for i, t := range expiredTrx {
		idsToExpire[i] = t.ID
	}

	if err := u.trxRepo.UpdateStatusInBatch(ctx, idsToExpire, entities.Expired); err != nil {
		return fmt.Errorf("failed to update status for expired payments: %w", err)
	}

	fmt.Printf("Expired %d payment transactions.\n", len(idsToExpire))
	return nil
}
