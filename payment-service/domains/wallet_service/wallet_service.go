package walletservice

import "github.com/gin-gonic/gin"

type WalletServiceClient interface {
	// VerifyBalance asks the core-service if a user has enough funds.
	VerifyBalance(ctx *gin.Context, userID string, amount float64) error
}
