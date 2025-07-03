package usecases

import (
	"login-api/domains/wallet"
	"login-api/domains/wallet/models/responses"

	"github.com/gin-gonic/gin"
)

type WalletUseCaseImpl struct {
	walletRepo wallet.WalletRepository
}

func NewWalletUsecase(walletRepo wallet.WalletRepository) *WalletUseCaseImpl {
	return &WalletUseCaseImpl{walletRepo: walletRepo}
}

func (u *WalletUseCaseImpl) GetBalance(ctx *gin.Context, userID string) (*responses.WalletBalanceResponse, error) {
	wallet, err := u.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err // Propagate the error (e.g., record not found).
	}

	response := &responses.WalletBalanceResponse{
		UserID:  wallet.UserID,
		Balance: wallet.Balance,
	}

	return response, nil
}
