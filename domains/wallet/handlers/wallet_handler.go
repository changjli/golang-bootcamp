package handlers

import (
	"errors"
	"login-api/domains/wallet"
	"login-api/helpers"
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
