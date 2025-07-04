package handlers

import (
	"core-service/domains/wallet"
	"core-service/domains/wallet/models/requests"
	"core-service/helpers"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WalletHandler struct {
	walletUsecase wallet.WalletUsecase
}

func NewWalletHandler(walletUsecase wallet.WalletUsecase) *WalletHandler {
	return &WalletHandler{walletUsecase: walletUsecase}
}

func (h *WalletHandler) GetBalance(ctx *gin.Context) {
	userID, err := helpers.GetAuthenticatedUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	balanceResponse, err := h.walletUsecase.GetBalance(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found for this user"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve wallet balance"})
		return
	}

	ctx.JSON(http.StatusOK, balanceResponse)
}

func (h *WalletHandler) VerifyBalance(ctx *gin.Context) {
	var req requests.VerifyBalanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Get the user's current balance using the existing usecase.
	wallet, err := h.walletUsecase.GetBalance(ctx, req.UserID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user_not_found"})
		return
	}

	// Check for sufficient funds.
	if wallet.Balance < req.Amount {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "insufficient_funds"})
		return
	}

	// If funds are sufficient, return success.
	ctx.JSON(http.StatusOK, gin.H{"status": "sufficient"})
}
