package repositories

import (
	"core-service/domains/transaction/entities"
	"core-service/infrastructures"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionRepositoryImpl struct {
	db infrastructures.Database
}

func NewTransactionRepository(db infrastructures.Database) *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{db: db}
}

func (r *TransactionRepositoryImpl) CreateInTx(ctx *gin.Context, tx *gorm.DB, transaction *entities.Transaction) error {
	return tx.WithContext(ctx).Create(transaction).Error
}

func (r *TransactionRepositoryImpl) GetHistoryByUserID(ctx *gin.Context, userID int, page int, limit int) ([]entities.Transaction, int64, error) {
	var transactions []entities.Transaction
	var total int64
	offset := (page - 1) * limit

	err := r.db.GetInstance().WithContext(ctx).Model(&entities.Transaction{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.GetInstance().WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Limit(limit).Offset(offset).Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}
