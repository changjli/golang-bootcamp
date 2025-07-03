package usecases

import (
	"context"
	"errors"
	"fmt"
	"login-api/domains/transaction"
	"login-api/domains/transaction/entities"
	"login-api/domains/transaction/models/requests"
	"login-api/domains/transaction/models/responses"
	"login-api/domains/wallet"
	"login-api/infrastructures"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionUseCaseImpl struct {
	db         infrastructures.Database
	trxRepo    transaction.TransactionRepository
	walletRepo wallet.WalletRepository
}

func NewTransactionUsecase(db infrastructures.Database, trxRepo transaction.TransactionRepository, walletRepo wallet.WalletRepository) *TransactionUseCaseImpl {
	return &TransactionUseCaseImpl{
		db:         db,
		trxRepo:    trxRepo,
		walletRepo: walletRepo,
	}
}

func (u *TransactionUseCaseImpl) TopUp(ctx *gin.Context, userID string, req *requests.TopUpRequest) (*responses.TopUpResponse, error) {
	var response *responses.TopUpResponse
	err := u.db.GetInstance().Transaction(func(tx *gorm.DB) error {
		// Get the user's wallet and lock it for the update.
		currentWallet, err := u.walletRepo.GetByUserID(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to get wallet: %w", err)
		}

		// Calculate the new balance.
		newBalance := currentWallet.Balance + req.Amount

		// Create the transaction record.
		topUpTrx := &entities.Transaction{
			ID:     uuid.NewString(),
			UserID: userID,
			Type:   entities.TopUp,
			Amount: req.Amount,
			Status: entities.Completed,
		}
		if err := u.trxRepo.CreateInTx(ctx, tx, topUpTrx); err != nil {
			return fmt.Errorf("failed to create transaction record: %w", err)
		}

		// Update the wallet balance.
		if err := u.walletRepo.UpdateBalance(ctx, userID, newBalance); err != nil {
			return fmt.Errorf("failed to update wallet balance: %w", err)
		}

		response = &responses.TopUpResponse{
			TransactionID: topUpTrx.ID,
			Message:       "Top-up successful",
			NewBalance:    newBalance,
		}
		return nil
	})

	return response, err
}

func (u *TransactionUseCaseImpl) Transfer(ctx *gin.Context, fromUserID string, req *requests.TransferRequest) (*responses.TransferResponse, error) {
	// Business rule: a user cannot transfer to themselves.
	if fromUserID == req.ToUserID {
		return nil, errors.New("cannot transfer to the same account")
	}

	var response *responses.TransferResponse
	err := u.db.GetInstance().Transaction(func(tx *gorm.DB) error {
		// Get sender's wallet.
		senderWallet, err := u.walletRepo.GetByUserID(ctx, fromUserID)
		if err != nil {
			return fmt.Errorf("sender wallet not found: %w", err)
		}

		// Check for sufficient funds.
		if senderWallet.Balance < req.Amount {
			return errors.New("insufficient funds")
		}

		// Get receiver's wallet.
		receiverWallet, err := u.walletRepo.GetByUserID(ctx, req.ToUserID)
		if err != nil {
			return fmt.Errorf("receiver wallet not found: %w", err)
		}

		// Perform the transfer.
		senderNewBalance := senderWallet.Balance - req.Amount
		receiverNewBalance := receiverWallet.Balance + req.Amount

		// Create transaction record for the sender (transfer_out).
		transferOutTrx := &entities.Transaction{
			ID:         uuid.NewString(),
			UserID:     fromUserID,
			Type:       entities.TransferOut,
			Amount:     req.Amount,
			Status:     entities.Completed,
			FromUserID: fromUserID,
			ToUserID:   req.ToUserID,
		}
		if err := u.trxRepo.CreateInTx(ctx, tx, transferOutTrx); err != nil {
			return fmt.Errorf("failed to create sender transaction: %w", err)
		}

		// Create transaction record for the receiver (transfer_in).
		transferInTrx := &entities.Transaction{
			ID:         uuid.NewString(),
			UserID:     req.ToUserID,
			Type:       entities.TransferIn,
			Amount:     req.Amount,
			Status:     entities.Completed,
			FromUserID: fromUserID,
			ToUserID:   req.ToUserID,
		}
		if err := u.trxRepo.CreateInTx(ctx, tx, transferInTrx); err != nil {
			return fmt.Errorf("failed to create receiver transaction: %w", err)
		}

		// Update both wallet balances.
		if err := u.walletRepo.UpdateBalance(ctx, fromUserID, senderNewBalance); err != nil {
			return fmt.Errorf("failed to update sender balance: %w", err)
		}
		if err := u.walletRepo.UpdateBalance(ctx, req.ToUserID, receiverNewBalance); err != nil {
			return fmt.Errorf("failed to update receiver balance: %w", err)
		}

		response = &responses.TransferResponse{
			TransactionID: transferOutTrx.ID,
			Status:        string(entities.Completed),
			Message:       "Transfer successful",
		}
		return nil
	})

	return response, err
}

// InitiatePayment creates a pending payment transaction.
func (u *TransactionUseCaseImpl) InitiatePayment(ctx *gin.Context, userID string, req *requests.PayRequest) (*responses.PayResponse, error) {
	// Get current wallet to check balance and return it in the response.
	currentWallet, err := u.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("wallet not found: %w", err)
	}

	// Business rule: check if the user has enough funds to even initiate the payment.
	if currentWallet.Balance < req.Amount {
		return nil, errors.New("insufficient funds to initiate payment")
	}

	// Create the pending payment transaction record.
	// NOTE: We do NOT use a DB transaction here and do NOT update the balance.
	paymentTrx := &entities.Transaction{
		ID:          uuid.NewString(),
		UserID:      userID,
		Type:        entities.Payment,
		Amount:      req.Amount,
		Status:      entities.Pending, // The key part of the requirement.
		MerchantID:  req.MerchantID,
		Description: req.Description,
	}

	// We use a temporary transaction here just for creating the single record.
	err = u.db.GetInstance().Transaction(func(tx *gorm.DB) error {
		return u.trxRepo.CreateInTx(ctx, tx, paymentTrx)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}

	response := &responses.PayResponse{
		TransactionID: paymentTrx.ID,
		Message:       "Payment initiated successfully. Awaiting confirmation.",
		NewBalance:    currentWallet.Balance, // Return the current, unchanged balance.
	}
	return response, nil
}

// GetHistory retrieves a user's transaction history.
func (u *TransactionUseCaseImpl) GetHistory(ctx *gin.Context, userID string, page int, limit int) (*responses.TransactionHistoryResponse, error) {
	transactions, total, err := u.trxRepo.GetHistoryByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve transaction history: %w", err)
	}

	details := make([]responses.TransactionDetail, len(transactions))
	for i, t := range transactions {
		details[i] = responses.TransactionDetail{
			TransactionID: t.ID,
			Type:          string(t.Type),
			Amount:        t.Amount,
			From:          t.FromUserID,
			To:            t.ToUserID,
			Timestamp:     t.CreatedAt,
			Status:        string(t.Status),
		}
	}

	response := &responses.TransactionHistoryResponse{
		Transactions: details,
		Pagination: responses.Pagination{
			CurrentPage: page,
			TotalPages:  int(math.Ceil(float64(total) / float64(limit))),
			TotalItems:  total,
		},
	}
	return response, nil
}

// Run this using cron or scheduler
// ExpirePayments finds and expires old pending payments.
func (u *TransactionUseCaseImpl) ExpirePayments(ctx context.Context) error {
	// Find all pending payments older than 10 minutes.
	expirationTime := time.Now().Add(-10 * time.Minute)
	expiredTrx, err := u.trxRepo.FindPendingPaymentsBefore(ctx, expirationTime)
	if err != nil {
		return fmt.Errorf("failed to find expired payments: %w", err)
	}

	if len(expiredTrx) == 0 {
		return nil // Nothing to do.
	}

	// Collect the IDs to update them in a single batch query.
	idsToExpire := make([]string, len(expiredTrx))
	for i, t := range expiredTrx {
		idsToExpire[i] = t.ID
	}

	// Update their status to 'expired'.
	if err := u.trxRepo.UpdateStatusInBatch(ctx, idsToExpire, entities.Expired); err != nil {
		return fmt.Errorf("failed to update status for expired payments: %w", err)
	}

	fmt.Printf("Expired %d payment transactions.\n", len(idsToExpire))
	return nil
}
