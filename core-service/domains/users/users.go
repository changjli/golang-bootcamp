package users

import (
	"core-service/domains/users/entities"
	"core-service/domains/users/models"

	"github.com/gin-gonic/gin"
)

type UserUseCase interface {
	Register(ctx *gin.Context, request models.RegisterReq) (*models.RegisterResp, error)
	Login(ctx *gin.Context, request models.LoginReq) (*models.LoginResp, error)
	Me(ctx *gin.Context) (*models.MeResp, error)
	Logout(ctx *gin.Context) error
}

type UserRepository interface {
	Save(ctx *gin.Context, user *entities.User) (*entities.User, error)
	FindbyName(ctx *gin.Context, name string) (*entities.User, error)
	FindbyEmail(ctx *gin.Context, email string) (*entities.User, error)
	FindbyID(ctx *gin.Context, id int) (*entities.User, error)
}
