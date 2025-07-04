package repositories

import (
	"errors"
	"fmt"
	"payment-service/domains/users/entities"
	"payment-service/infrastructures"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var users []entities.User

type UserRepository struct {
	db infrastructures.Database
}

func NewUserRepository(db infrastructures.Database) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Save(ctx *gin.Context, user *entities.User) (*entities.User, error) {
	if err := r.db.GetInstance().WithContext(ctx).Create(user).Error; err != nil {
		return nil, fmt.Errorf("[USER_REPOSITORY]: failed to save user: %w", err)
	}

	return user, nil

}

func (r *UserRepository) FindByUsername(ctx *gin.Context, username string) (*entities.User, error) {
	var user entities.User
	if err := r.db.GetInstance().WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("[USER_REPOSITORY]: User not found")
		}
		return nil, fmt.Errorf("[USER_REPOSITORY]: failed to find user: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) FindByUsernameAndPassword(ctx *gin.Context, username string, password string) (*entities.User, error) {
	// First, find the user by their username.
	user, err := r.FindByUsername(ctx, username)
	if err != nil {
		// Return a generic error to avoid revealing whether the username exists.
		return nil, fmt.Errorf("[USER_REPOSITORY]: Invalid credentials")
	}

	// Compare the provided password with the stored hashed password.
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// Passwords do not match.
		return nil, fmt.Errorf("[USER_REPOSITORY]: Invalid credentials")
	}

	// Password is correct.
	return user, nil
}
