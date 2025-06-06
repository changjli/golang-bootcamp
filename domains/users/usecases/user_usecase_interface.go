package usecases

import (
	"login-api/domains/users/models/requests"
	"login-api/domains/users/models/responses"

	"github.com/gin-gonic/gin"
)

type UserUseCaseInterface interface {
	Login(ctx *gin.Context, request requests.UserLoginRequest) (*responses.UserLoginResponse, error)
	Me(ctx *gin.Context) (*responses.UserMeResponse, error)
}
