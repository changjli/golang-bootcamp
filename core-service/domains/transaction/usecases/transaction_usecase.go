package usecases

import (
	"core-service/domains/transaction"
	"core-service/domains/transaction/entities"
	"core-service/domains/transaction/models/requests"
	"core-service/domains/transaction/models/responses"
	"core-service/domains/wallet"
	"core-service/infrastructures"
	"core-service/infrastructures/messaging"
	"errors"
	"events"
	"fmt"
	"math"
	"time"
	"utils/helpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionUseCaseImpl struct {
	db         infrastructures.Database
	trxRepo    transaction.TransactionRepository
	walletRepo wallet.WalletRepository
	publisher  messaging.PaymentPublisher
}

func NewTransactionUsecase(db infrastructures.Database, trxRepo transaction.TransactionRepository, walletRepo wallet.WalletRepository, publisher messaging.PaymentPublisher) *TransactionUseCaseImpl {
	return &TransactionUseCaseImpl{
		db:         db,
		trxRepo:    trxRepo,
		walletRepo: walletRepo,
		publisher:  publisher,
	}
}

func (u *TransactionUseCaseImpl) TopUp(ctx *gin.Context, userID int, req *requests.TopUpRequest) (*responses.TopUpResponse, error) {
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

func (u *TransactionUseCaseImpl) Transfer(ctx *gin.Context, fromUserID int, req *requests.TransferRequest) (*responses.TransferResponse, error) {
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

// GetHistory retrieves a user's transaction history.
func (u *TransactionUseCaseImpl) GetHistory(ctx *gin.Context, userID int, page int, limit int) (*responses.TransactionHistoryResponse, error) {
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

// InitiatePayment now publishes an event instead of creating a DB record.
func (u *TransactionUseCaseImpl) InitiatePayment(ctx *gin.Context, userID int, req *requests.PayRequest) (*responses.PayResponse, error) {
	claims, err := helpers.GetAuthenticatedClaims(ctx)
	if err != nil {
		return nil, err
	}

	// 1. Check if the user has sufficient funds. This validation remains in the core-service.
	currentWallet, err := u.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("wallet not found: %w", err)
	}
	if currentWallet.Balance < req.Amount {
		return nil, errors.New("insufficient funds")
	}

	// 3. Create the event payload.
	paymentID := uuid.NewString()
	event := &events.PaymentRequestedEvent{
		PaymentID:   paymentID,
		UserID:      userID,
		Amount:      req.Amount,
		CreatedAt:   time.Now(),
		AuthToken:   claims.TokenString,
		Description: req.Description,
		MerchantID:  req.MerchantID,
	}

	// 4. Publish the event to RabbitMQ.
	if err := u.publisher.PublishPaymentRequested(ctx, event); err != nil {
		// If publishing fails, the payment cannot be initiated.
		return nil, fmt.Errorf("failed to publish payment request: %w", err)
	}

	// 5. Immediately return a "pending" response to the user.
	response := &responses.PayResponse{
		TransactionID: paymentID,
		Message:       "Payment request received and is being processed.",
		NewBalance:    currentWallet.Balance, // The balance is not yet changed.
	}
	return response, nil
}
