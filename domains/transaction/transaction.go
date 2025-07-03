package transaction

import (
	"context"
	"login-api/domains/transaction/entities"
	"login-api/domains/transaction/models/requests"
	"login-api/domains/transaction/models/responses"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TransactionRepository defines the database operations for transactions.
type TransactionRepository interface {
	// CreateInTx creates a new transaction record within a provided database transaction.
	// This ensures that transaction creation can be part of a larger atomic operation
	// (e.g., one that also updates a wallet balance).
	CreateInTx(ctx *gin.Context, tx *gorm.DB, transaction *entities.Transaction) error

	// GetHistoryByUserID retrieves a paginated list of transactions for a specific user.
	// It also returns the total count of transactions for calculating total pages.
	GetHistoryByUserID(ctx *gin.Context, userID string, page int, limit int) ([]entities.Transaction, int64, error)

	// FindPendingPaymentsBefore retrieves all payment transactions that are in 'pending' status
	// and were created before the specified expiration time. This is used by the background
	// worker to find payments that need to be marked as 'expired'.
	FindPendingPaymentsBefore(ctx context.Context, expirationTime time.Time) ([]entities.Transaction, error)

	// UpdateStatusInBatch updates the status of multiple transactions at once using their IDs.
	// This is more efficient than updating them one by one.
	UpdateStatusInBatch(ctx context.Context, transactionIDs []string, status entities.TransactionStatus) error
}

// TransactionUsecase defines the business logic operations for transactions.
type TransactionUsecase interface {
	// TopUp handles the logic for a user adding funds to their wallet.
	// This involves creating a transaction record and updating the wallet balance atomically.
	TopUp(ctx *gin.Context, userID string, req *requests.TopUpRequest) (*responses.TopUpResponse, error)

	// Transfer handles the logic for a user sending funds to another user.
	// This is a complex operation that must atomically decrease the sender's balance,
	// increase the receiver's balance, and create two transaction records (one 'transfer_out'
	// and one 'transfer_in').
	Transfer(ctx *gin.Context, fromUserID string, req *requests.TransferRequest) (*responses.TransferResponse, error)

	// InitiatePayment handles the special requirement for creating an expiring payment session.
	// It creates a 'pending' payment transaction but does NOT deduct the user's balance.
	InitiatePayment(ctx *gin.Context, userID string, req *requests.PayRequest) (*responses.PayResponse, error)

	// GetHistory retrieves a user's transaction history with pagination.
	GetHistory(ctx *gin.Context, userID string, page int, limit int) (*responses.TransactionHistoryResponse, error)

	// ExpirePayments is the method called by a background worker. It finds all pending
	// payments that are older than the 10-minute window and updates their status to 'expired'.
	// It uses a standard context.Context as it is not tied to a specific API request.
	ExpirePayments(ctx context.Context) error
}
