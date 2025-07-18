package repositories

import (
	"core-service/domains/users/entities"
	"core-service/infrastructures"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var users []entities.User

type UserRepositoryImpl struct {
	db infrastructures.Database
}

func NewUserRepository(db infrastructures.Database) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) Save(ctx *gin.Context, user *entities.User) (*entities.User, error) {
	if err := r.db.GetInstance().WithContext(ctx).Create(user).Error; err != nil {
		return nil, fmt.Errorf("[USER_REPOSITORY]: failed to save user: %w", err)
	}

	return user, nil

}

func (r *UserRepositoryImpl) FindbyName(ctx *gin.Context, name string) (*entities.User, error) {
	var user entities.User
	if err := r.db.GetInstance().WithContext(ctx).Where("name = ?", name).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user is not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindbyEmail(ctx *gin.Context, email string) (*entities.User, error) {
	var user entities.User
	if err := r.db.GetInstance().WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user is not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindbyID(ctx *gin.Context, id int) (*entities.User, error) {
	var user entities.User
	if err := r.db.GetInstance().WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("[USER_REPOSITORY]: User not found")
		}
		return nil, fmt.Errorf("[USER_REPOSITORY]: failed to find user: %w", err)
	}
	return &user, nil
}
