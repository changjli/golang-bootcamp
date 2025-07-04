package wallet

import (
	"core-service/domains/wallet/entities"
	"core-service/domains/wallet/models/responses"

	"github.com/gin-gonic/gin"
)

type WalletRepository interface {
	// Create creates a new wallet for a given user ID, typically with an initial balance.
	Create(ctx *gin.Context, userID string) (*entities.Wallet, error)

	// GetByUserID retrieves a user's wallet using their user ID.
	GetByUserID(ctx *gin.Context, userID string) (*entities.Wallet, error)

	// UpdateBalance updates the balance of a wallet identified by its user ID.
	// This operation should be atomic. It's often used within a database transaction
	// alongside creating a transaction record.
	UpdateBalance(ctx *gin.Context, userID string, newBalance float64) error
}

type WalletUsecase interface {
	// GetBalance retrieves the current balance for a specific user.
	// It takes a user ID and returns a DTO suitable for the API response.
	GetBalance(ctx *gin.Context, userID string) (*responses.WalletBalanceResponse, error)
}
