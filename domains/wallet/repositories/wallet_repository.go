package repositories

import (
	"errors"
	"login-api/domains/wallet/entities"
	"login-api/infrastructures"

	"github.com/gin-gonic/gin"
)

type WalletRepositoryImpl struct {
	db infrastructures.Database
}

func NewWalletRepository(db infrastructures.Database) *WalletRepositoryImpl {
	return &WalletRepositoryImpl{db: db}
}

func (r *WalletRepositoryImpl) Create(ctx *gin.Context, userID string) (*entities.Wallet, error) {
	wallet := &entities.Wallet{
		UserID:  userID,
		Balance: 0,
	}

	if err := r.db.GetInstance().WithContext(ctx).Create(wallet).Error; err != nil {
		return nil, err
	}
	return wallet, nil
}

func (r *WalletRepositoryImpl) GetByUserID(ctx *gin.Context, userID string) (*entities.Wallet, error) {
	var wallet entities.Wallet
	if err := r.db.GetInstance().WithContext(ctx).Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepositoryImpl) UpdateBalance(ctx *gin.Context, userID string, newBalance float64) error {
	result := r.db.GetInstance().WithContext(ctx).Model(&entities.Wallet{}).Where("user_id = ?", userID).Update("balance", newBalance)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("wallet not found or no rows affected")
	}
	return nil
}
