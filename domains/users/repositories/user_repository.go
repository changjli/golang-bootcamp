package repositories

import (
	"fmt"
	"login-api/domains/users/entities"
	"login-api/helpers"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var users []entities.User

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	users = append(users, entities.User{
		Id:       1,
		Username: "Lui",
		Name:     "Luigi",
		Password: helpers.GenerateHashedPassword("12345678"),
	})
	return &UserRepository{}
}

func (r *UserRepository) Save(ctx *gin.Context, user entities.User) *entities.User {
	panic("not implemented") // TODO: Implement
}

func (r *UserRepository) FindByUsername(ctx *gin.Context, username string) (*entities.User, error) {
	for _, user := range users {
		if user.Username == username {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("[USER_REPOSITORY]: User not found")
}

func (r *UserRepository) FindByUsernameAndPassword(ctx *gin.Context, username string, password string) (*entities.User, error) {
	for _, user := range users {
		if user.Username == username {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err == nil {
				return &user, nil
			}
		}
	}

	return nil, fmt.Errorf("[USER_REPOSITORY]: User not found")
}
