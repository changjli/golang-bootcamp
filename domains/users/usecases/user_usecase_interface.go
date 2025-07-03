package usecases

import (
	"login-api/domains/users/models/requests"
	"login-api/domains/users/models/responses"

	"github.com/gin-gonic/gin"
)

type UserUseCaseInterface interface {
	Register(ctx *gin.Context, request requests.UserRegisterRequest) (*responses.UserRegisterResponse, error)
	Login(ctx *gin.Context, request requests.UserLoginRequest) (*responses.UserLoginResponse, error)
	Me(ctx *gin.Context) (*responses.UserMeResponse, error)
	Logout(ctx *gin.Context) error
}
