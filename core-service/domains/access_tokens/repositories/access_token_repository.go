package repositories

import (
	"errors"
	"core-service/domains/access_tokens/entities"
	"core-service/infrastructures"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("record not found")

type AccessTokenRepository struct {
	db infrastructures.Database
}

func NewAccessTokenRepository(db infrastructures.Database) *AccessTokenRepository {
	return &AccessTokenRepository{
		db: db,
	}
}

func (r *AccessTokenRepository) Save(ctx *gin.Context, token *entities.AccessToken) error {
	tx := r.db.GetInstance().WithContext(ctx).Create(token)
	return tx.Error
}

func (r *AccessTokenRepository) FindByID(ctx *gin.Context, id string) (*entities.AccessToken, error) {
	var token entities.AccessToken
	tx := r.db.GetInstance().WithContext(ctx).Where("id = ?", id).First(&token)

	if tx.Error != nil {
		// Check if the error is a "record not found" error and return our custom error.
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		// Return any other database error.
		return nil, tx.Error
	}

	return &token, nil
}

func (r *AccessTokenRepository) Revoke(ctx *gin.Context, id string) error {
	tx := r.db.GetInstance().WithContext(ctx).Model(&entities.AccessToken{}).Where("id = ?", id).Update("revoked", true)

	if tx.Error != nil {
		return tx.Error
	}

	// If no rows were affected, it means the token ID did not exist.
	if tx.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
