package repositories

import (
	"login-api/domains/users/entities"

	"github.com/gin-gonic/gin"
)

type UserRepositoryInterface interface {
	Save(ctx *gin.Context, user *entities.User) (*entities.User, error)
	FindByUsername(ctx *gin.Context, username string) (*entities.User, error)
	FindByUsernameAndPassword(ctx *gin.Context, username string, password string) (*entities.User, error)
}
