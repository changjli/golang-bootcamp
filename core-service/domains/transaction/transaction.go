package transaction

import (
	"core-service/domains/transaction/entities"
	"core-service/domains/transaction/models/requests"
	"core-service/domains/transaction/models/responses"

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
	GetHistoryByUserID(ctx *gin.Context, userID int, page int, limit int) ([]entities.Transaction, int64, error)
}

// TransactionUsecase defines the business logic operations for transactions.
type TransactionUsecase interface {
	// TopUp handles the logic for a user adding funds to their wallet.
	// This involves creating a transaction record and updating the wallet balance atomically.
	TopUp(ctx *gin.Context, userID int, req *requests.TopUpRequest) (*responses.TopUpResponse, error)

	// Transfer handles the logic for a user sending funds to another user.
	// This is a complex operation that must atomically decrease the sender's balance,
	// increase the receiver's balance, and create two transaction records (one 'transfer_out'
	// and one 'transfer_in').
	Transfer(ctx *gin.Context, fromUserID int, req *requests.TransferRequest) (*responses.TransferResponse, error)

	// GetHistory retrieves a user's transaction history with pagination.
	GetHistory(ctx *gin.Context, userID int, page int, limit int) (*responses.TransactionHistoryResponse, error)
}
